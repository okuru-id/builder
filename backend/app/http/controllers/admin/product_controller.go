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
	facades.Orm().Query().OrderBy("created_at", "desc").Get(&products)
	return ctx.Response().Success().Json(http.Json{"data": products})
}

func (c *ProductController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", id).First(&product); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": product})
}

func (c *ProductController) Store(ctx http.Context) http.Response {
	var product models.Product
	if err := ctx.Request().Bind(&product); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if err := facades.Orm().Query().Create(&product); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": product})
}

func (c *ProductController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", id).First(&product); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&product); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	facades.Orm().Query().Save(&product)
	return ctx.Response().Success().Json(http.Json{"data": product})
}

func (c *ProductController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	facades.Orm().Query().Where("id = ?", id).Delete(&models.Product{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
