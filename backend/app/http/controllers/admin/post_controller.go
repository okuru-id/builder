package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (c *PostController) Index(ctx http.Context) http.Response {
	var posts []models.Post
	facades.Orm().Query().OrderBy("created_at", "desc").Get(&posts)
	return ctx.Response().Success().Json(http.Json{"data": posts})
}

func (c *PostController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var post models.Post
	if err := facades.Orm().Query().Where("id = ?", id).First(&post); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": post})
}

func (c *PostController) Store(ctx http.Context) http.Response {
	var post models.Post
	if err := ctx.Request().Bind(&post); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if err := facades.Orm().Query().Create(&post); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": post})
}

func (c *PostController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var post models.Post
	if err := facades.Orm().Query().Where("id = ?", id).First(&post); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&post); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	facades.Orm().Query().Save(&post)
	return ctx.Response().Success().Json(http.Json{"data": post})
}

func (c *PostController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	facades.Orm().Query().Where("id = ?", id).Delete(&models.Post{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
