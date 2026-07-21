package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type ProductController struct{}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (c *ProductController) Index(ctx http.Context) http.Response {
	var products []models.Product
	facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).OrderBy("created_at", "desc").Get(&products)
	return ctx.Response().Success().Json(http.Json{"data": products})
}

func (c *ProductController) Show(ctx http.Context) http.Response {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&product); err != nil || product.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": product})
}

func (c *ProductController) Store(ctx http.Context) http.Response {
	var product models.Product
	if err := ctx.Request().Bind(&product); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	product.UserID = currentUserID(ctx)
	if err := facades.Orm().Query().Create(&product); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": product})
}

func (c *ProductController) Update(ctx http.Context) http.Response {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&product); err != nil || product.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&product); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	product.UserID = currentUserID(ctx)
	facades.Orm().Query().Save(&product)
	return ctx.Response().Success().Json(http.Json{"data": product})
}

func (c *ProductController) Destroy(ctx http.Context) http.Response {
	facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).Delete(&models.Product{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
