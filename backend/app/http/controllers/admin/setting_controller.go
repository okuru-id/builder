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

// Index lists all settings as a key-value map.
func (c *SettingController) Index(ctx http.Context) http.Response {
	var settings []models.Setting
	facades.Orm().Query().OrderBy("key", "asc").Get(&settings)

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

// Update creates or updates a setting by key.
func (c *SettingController) Update(ctx http.Context) http.Response {
	var input updateSettingRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if input.Key == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "key is required"})
	}

	var setting models.Setting
	err := facades.Orm().Query().Where("key = ?", input.Key).First(&setting)
	if err != nil {
		setting = models.Setting{Key: input.Key, Value: input.Value}
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
	key := ctx.Request().Route("key")
	facades.Orm().Query().Where("key = ?", key).Delete(&models.Setting{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
