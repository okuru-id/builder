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
	facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).OrderBy("created_at", "desc").Get(&messages)
	return ctx.Response().Success().Json(http.Json{"data": messages})
}

func (c *MessageController) Show(ctx http.Context) http.Response {
	var message models.Message
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&message); err != nil || message.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": message})
}

func (c *MessageController) MarkRead(ctx http.Context) http.Response {
	_, err := facades.Orm().Query().Model(&models.Message{}).
		Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).
		Update("status", "read")
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"updated": true})
}

func (c *MessageController) Archive(ctx http.Context) http.Response {
	_, err := facades.Orm().Query().Model(&models.Message{}).
		Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).
		Update("status", "archived")
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"updated": true})
}

func (c *MessageController) Destroy(ctx http.Context) http.Response {
	facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).Delete(&models.Message{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
