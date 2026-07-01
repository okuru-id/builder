package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type MessageController struct{}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (c *MessageController) Index(ctx http.Context) http.Response {
	var messages []models.Message
	facades.Orm().Query().OrderBy("created_at", "desc").Get(&messages)
	return ctx.Response().Success().Json(http.Json{"data": messages})
}

func (c *MessageController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var message models.Message
	if err := facades.Orm().Query().Where("id = ?", id).First(&message); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": message})
}

func (c *MessageController) MarkRead(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if _, err := facades.Orm().Query().Model(&models.Message{}).
		Where("id = ?", id).
		Update("status", "read"); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"updated": true})
}

func (c *MessageController) Archive(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if _, err := facades.Orm().Query().Model(&models.Message{}).
		Where("id = ?", id).
		Update("status", "archived"); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"updated": true})
}

func (c *MessageController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	facades.Orm().Query().Where("id = ?", id).Delete(&models.Message{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
