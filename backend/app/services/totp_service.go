package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"image/png"
	"io"
	"strings"

	"github.com/pquerna/otp/totp"

	"okuru/app/facades"
)

// TotpEnabled reports whether 2FA/TOTP is active.
// Explicit TOTP_ENABLED env overrides; otherwise auto-disabled
// for non-production envs (local, dev, testing).
func TotpEnabled() bool {
	env := strings.ToLower(facades.Config().GetString("app.env", "production"))
	def := env != "local" && env != "dev" && env != "testing"
	return facades.Config().EnvBool("TOTP_ENABLED", def)
}

type TotpService struct {
	appKey []byte
}

func NewTotpService(appKey string) *TotpService {
	key := []byte(appKey)
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		key = padded
	}
	return &TotpService{appKey: key[:32]}
}

func (s *TotpService) GenerateSecret(email string) (secret string, qrUrl string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Okuru.id",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	image, err := key.Image(176, 176)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, image); err != nil {
		return "", "", err
	}

	return key.Secret(), "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (s *TotpService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (s *TotpService) EncryptSecret(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.appKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base32.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *TotpService) DecryptSecret(encoded string) (string, error) {
	ciphertext, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.appKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
