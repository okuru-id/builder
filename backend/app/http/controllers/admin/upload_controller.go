package admin

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) Store(ctx http.Context) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "No file uploaded"})
	}

	size, err := file.Size()
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "Failed to read file size"})
	}
	if size > 5*1024*1024 {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "File too large (max 5MB)"})
	}

	ext := strings.ToLower(file.GetClientOriginalExtension())
	allowed := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true,
		"webp": true, "pdf": true,
	}
	if !allowed[ext] {
		// Fall back to parsing extension from original name.
		ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(file.GetClientOriginalName()), "."))
		if !allowed[ext] {
			return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "File type not allowed"})
		}
	}

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	filename := hex.EncodeToString(randomBytes) + "." + ext

	path, err := file.StoreAs("uploads", filename)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": "Failed to store file"})
	}

	return ctx.Response().Success().Json(http.Json{
		"path": path,
		"url":  "/storage/" + path,
	})
}
