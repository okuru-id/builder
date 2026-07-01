package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type OpenSourceController struct{}

func NewOpenSourceController() *OpenSourceController {
	return &OpenSourceController{}
}

func (c *OpenSourceController) Index(ctx http.Context) http.Response {
	var projects []models.OpenSourceProject
	facades.Orm().Query().OrderBy("sort_order", "asc").Get(&projects)
	return ctx.Response().Success().Json(http.Json{"data": projects})
}

func (c *OpenSourceController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var project models.OpenSourceProject
	if err := facades.Orm().Query().Where("id = ?", id).First(&project); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": project})
}

func (c *OpenSourceController) Store(ctx http.Context) http.Response {
	var project models.OpenSourceProject
	if err := ctx.Request().Bind(&project); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if err := facades.Orm().Query().Create(&project); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": project})
}

func (c *OpenSourceController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var project models.OpenSourceProject
	if err := facades.Orm().Query().Where("id = ?", id).First(&project); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&project); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	facades.Orm().Query().Save(&project)
	return ctx.Response().Success().Json(http.Json{"data": project})
}

func (c *OpenSourceController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	facades.Orm().Query().Where("id = ?", id).Delete(&models.OpenSourceProject{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
