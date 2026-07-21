package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type SettingController struct{}

func NewSettingController() *SettingController {
	return &SettingController{}
}

// Index lists the caller's own settings as a key-value map.
func (c *SettingController) Index(ctx http.Context) http.Response {
	var settings []models.Setting
	facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).OrderBy("key", "asc").Get(&settings)

	out := make(map[string]string, len(settings))
	for _, s := range settings {
		out[s.Key] = s.Value
	}
	return ctx.Response().Success().Json(http.Json{"data": out})
}

type updateSettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Update creates or updates a setting scoped to the caller.
func (c *SettingController) Update(ctx http.Context) http.Response {
	var input updateSettingRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if input.Key == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "key is required"})
	}

	uid := currentUserID(ctx)
	var setting models.Setting
	err := facades.Orm().Query().Where("key = ? AND user_id = ?", input.Key, uid).First(&setting)
	if err != nil || setting.ID == 0 {
		setting = models.Setting{Key: input.Key, Value: input.Value, UserID: uid}
		if createErr := facades.Orm().Query().Create(&setting); createErr != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": createErr.Error()})
		}
	} else {
		setting.Value = input.Value
		facades.Orm().Query().Save(&setting)
	}

	return ctx.Response().Success().Json(http.Json{"data": setting})
}

func (c *SettingController) Destroy(ctx http.Context) http.Response {
	facades.Orm().Query().Where("key = ? AND user_id = ?", ctx.Request().Route("key"), currentUserID(ctx)).Delete(&models.Setting{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
