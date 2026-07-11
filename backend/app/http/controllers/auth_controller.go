package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
	"okuru/app/services"
)

type AuthController struct {
	totp *services.TotpService
}

func NewAuthController() *AuthController {
	appKey := facades.Config().GetString("app.key", "")
	return &AuthController{
		totp: services.NewTotpService(appKey),
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login(ctx http.Context) http.Response {
	var input loginRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "invalid request body",
		})
	}

	if input.Email == "" || input.Password == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "email and password are required",
		})
	}

	var user models.User
	if err := facades.Orm().Query().Where("email = ?", input.Email).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "invalid credentials",
		})
	}

	if !facades.Hash().Check(input.Password, user.Password) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "invalid credentials",
		})
	}

	if services.TotpEnabled() && user.TotpSecret != nil && user.TotpVerified {
		tempToken, err := facades.Auth(ctx).LoginUsingID(user.ID)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"error": "failed to issue token",
			})
		}
		return ctx.Response().Json(http.StatusOK, http.Json{
			"requires_totp": true,
			"temp_token":    tempToken,
		})
	}

	token, err := facades.Auth(ctx).LoginUsingID(user.ID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to issue token",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"requires_totp":       false,
		"access_token":        token,
		"token_type":          "bearer",
		"totp_setup_required": services.TotpEnabled() && (user.TotpSecret == nil || !user.TotpVerified),
	})
}

type verifyTotpRequest struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
}

func (c *AuthController) VerifyTotp(ctx http.Context) http.Response {
	var input verifyTotpRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "invalid request body",
		})
	}

	if input.TempToken == "" || input.Code == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "temp_token and code are required",
		})
	}

	token := strings.TrimPrefix(input.TempToken, "Bearer ")
	payload, err := facades.Auth(ctx).Parse(token)
	if err != nil || payload == nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "invalid or expired temp token",
		})
	}

	userID, _ := strconv.ParseUint(payload.Key, 10, 64)

	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "user not found",
		})
	}

	if user.TotpSecret == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "TOTP is not configured for this account",
		})
	}

	secret, err := c.totp.DecryptSecret(*user.TotpSecret)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to read TOTP secret",
		})
	}

	if !c.totp.ValidateCode(secret, input.Code) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "invalid TOTP code",
		})
	}

	accessToken, err := facades.Auth(ctx).LoginUsingID(user.ID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to issue token",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"access_token": accessToken,
		"token_type":   "bearer",
	})
}

func (c *AuthController) SetupTotp(ctx http.Context) http.Response {
	if !services.TotpEnabled() {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"error": "TOTP is disabled",
		})
	}

	userIDStr, err := facades.Auth(ctx).ID()
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "unauthorized",
		})
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "user not found",
		})
	}

	secret, qrUrl, err := c.totp.GenerateSecret(user.Email)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to generate TOTP secret",
		})
	}

	encrypted, err := c.totp.EncryptSecret(secret)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to store TOTP secret",
		})
	}

	if _, err := facades.Orm().Query().Model(&models.User{}).
		Where("id = ?", userID).
		Update(map[string]any{"totp_secret": encrypted, "totp_verified": false}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to save TOTP secret",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"secret":  secret,
		"qr_url":  qrUrl,
		"issuer":  "Okuru.id",
		"account": user.Email,
	})
}

type verifySetupRequest struct {
	Code string `json:"code"`
}

func (c *AuthController) VerifyTotpSetup(ctx http.Context) http.Response {
	var input verifySetupRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "invalid request body",
		})
	}

	if input.Code == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "code is required",
		})
	}

	userIDStr, err := facades.Auth(ctx).ID()
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "unauthorized",
		})
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "user not found",
		})
	}

	if user.TotpSecret == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "TOTP setup has not been initiated",
		})
	}

	secret, err := c.totp.DecryptSecret(*user.TotpSecret)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to read TOTP secret",
		})
	}

	if !c.totp.ValidateCode(secret, input.Code) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "invalid TOTP code",
		})
	}

	if _, err := facades.Orm().Query().Model(&models.User{}).
		Where("id = ?", userID).
		Update("totp_verified", true); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to verify TOTP",
		})
	}

	accessToken, err := facades.Auth(ctx).LoginUsingID(user.ID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "failed to issue token",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"verified":     true,
		"access_token": accessToken,
		"token_type":   "bearer",
	})
}

func (c *AuthController) Me(ctx http.Context) http.Response {
	userIDStr, err := facades.Auth(ctx).ID()
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "unauthorized",
		})
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "user not found",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"id":            user.ID,
		"email":         user.Email,
		"totp_verified": user.TotpVerified,
		"totp_enabled":  user.TotpSecret != nil,
	})
}
