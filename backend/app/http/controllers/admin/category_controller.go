package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type CategoryController struct{}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

func (c *CategoryController) Index(ctx http.Context) http.Response {
	var categories []models.Category
	facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).OrderBy("id", "asc").Get(&categories)
	return ctx.Response().Success().Json(http.Json{"data": categories})
}

func (c *CategoryController) Show(ctx http.Context) http.Response {
	var category models.Category
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&category); err != nil || category.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": category})
}

func (c *CategoryController) Store(ctx http.Context) http.Response {
	var category models.Category
	if err := ctx.Request().Bind(&category); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	category.UserID = currentUserID(ctx)
	if err := facades.Orm().Query().Create(&category); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": category})
}

func (c *CategoryController) Update(ctx http.Context) http.Response {
	var category models.Category
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&category); err != nil || category.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&category); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	category.UserID = currentUserID(ctx)
	facades.Orm().Query().Save(&category)
	return ctx.Response().Success().Json(http.Json{"data": category})
}

func (c *CategoryController) Destroy(ctx http.Context) http.Response {
	facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).Delete(&models.Category{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
