package admin

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (c *ProjectController) Index(ctx http.Context) http.Response {
	var projects []models.Project
	facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).OrderBy("sort_order", "asc").Get(&projects)
	return ctx.Response().Success().Json(http.Json{"data": projects})
}

func (c *ProjectController) Show(ctx http.Context) http.Response {
	var project models.Project
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&project); err != nil || project.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": project})
}

func (c *ProjectController) Store(ctx http.Context) http.Response {
	var project models.Project
	if err := ctx.Request().Bind(&project); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	project.UserID = currentUserID(ctx)
	if err := facades.Orm().Query().Create(&project); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": project})
}

func (c *ProjectController) Update(ctx http.Context) http.Response {
	var project models.Project
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).First(&project); err != nil || project.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	if err := ctx.Request().Bind(&project); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	project.UserID = currentUserID(ctx)
	facades.Orm().Query().Save(&project)
	return ctx.Response().Success().Json(http.Json{"data": project})
}

func (c *ProjectController) Destroy(ctx http.Context) http.Response {
	facades.Orm().Query().Where("id = ? AND user_id = ?", ctx.Request().Route("id"), currentUserID(ctx)).Delete(&models.Project{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
