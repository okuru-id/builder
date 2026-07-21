package admin

import (
	"encoding/json"
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"gorm.io/datatypes"

	"okuru/app/facades"
	"okuru/app/models"
)

// LandingComponentController manages reusable component masters, scoped per user.
type LandingComponentController struct{}

func NewLandingComponentController() *LandingComponentController {
	return &LandingComponentController{}
}

type landingComponentInput struct {
	Name string          `json:"name" validate:"required"`
	Tree json.RawMessage `json:"tree"`
}

func ownComponent(ctx http.Context, id string) (models.LandingComponent, http.Response, bool) {
	var comp models.LandingComponent
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", id, currentUserID(ctx)).First(&comp); err != nil || comp.ID == 0 {
		return comp, ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"}), false
	}
	return comp, nil, true
}

func (c *LandingComponentController) Index(ctx http.Context) http.Response {
	var comps []models.LandingComponent
	if err := facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).Order("id desc").Get(&comps); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comps})
}

func (c *LandingComponentController) Show(ctx http.Context) http.Response {
	comp, resp, ok := ownComponent(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

func (c *LandingComponentController) Store(ctx http.Context) http.Response {
	var in landingComponentInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	comp := models.LandingComponent{
		Name:     in.Name,
		Tree:     datatypes.JSON(in.Tree),
		IsSystem: false,
		UserID:   currentUserID(ctx),
	}
	if err := facades.Orm().Query().Create(&comp); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

func (c *LandingComponentController) Duplicate(ctx http.Context) http.Response {
	comp, resp, ok := ownComponent(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	copy := models.LandingComponent{
		Name:     fmt.Sprintf("%s (copy)", comp.Name),
		Tree:     comp.Tree,
		IsSystem: false,
		UserID:   currentUserID(ctx),
	}
	if err := facades.Orm().Query().Create(&copy); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": copy})
}

func (c *LandingComponentController) Update(ctx http.Context) http.Response {
	comp, resp, ok := ownComponent(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	var in landingComponentInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if comp.IsSystem {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "system component cannot be modified"})
	}
	if in.Name != "" {
		comp.Name = in.Name
	}
	if len(in.Tree) > 0 {
		comp.Tree = datatypes.JSON(in.Tree)
	}
	if _, err := facades.Orm().Query().Where("id = ?", comp.ID).Update(&comp); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

func (c *LandingComponentController) Destroy(ctx http.Context) http.Response {
	comp, resp, ok := ownComponent(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	if comp.IsSystem {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "system component cannot be deleted"})
	}
	if _, err := facades.Orm().Query().Where("id = ?", comp.ID).Delete(&models.LandingComponent{}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": true})
}
