package admin

import (
	"encoding/json"

	"github.com/goravel/framework/contracts/http"
	"gorm.io/datatypes"

	"okuru/app/facades"
	"okuru/app/models"
)

// LandingComponentController manages reusable component masters.
// Instances in page trees reference a master by ID; rendering resolves the
// master tree at render time so editing a master updates all instances.
type LandingComponentController struct{}

func NewLandingComponentController() *LandingComponentController {
	return &LandingComponentController{}
}

type landingComponentInput struct {
	Name string          `json:"name" validate:"required"`
	Tree json.RawMessage `json:"tree"`
}

// Index lists all component masters.
func (c *LandingComponentController) Index(ctx http.Context) http.Response {
	var comps []models.LandingComponent
	if err := facades.Orm().Query().Order("id desc").Get(&comps); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comps})
}

// Show returns a single component master with its tree.
func (c *LandingComponentController) Show(ctx http.Context) http.Response {
	var comp models.LandingComponent
	if err := facades.Orm().Query().Where("id = ?", ctx.Request().Input("id")).First(&comp); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	if comp.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

// Store creates a new component master from a tree fragment.
func (c *LandingComponentController) Store(ctx http.Context) http.Response {
	var in landingComponentInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	comp := models.LandingComponent{
		Name: in.Name,
		Tree: datatypes.JSON(in.Tree),
	}
	if err := facades.Orm().Query().Create(&comp); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

// Update replaces the master name and/or tree. Instances resolve the new tree
// on next render — this is how "edit master → update all instances" works.
func (c *LandingComponentController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var comp models.LandingComponent
	if err := facades.Orm().Query().Where("id = ?", id).First(&comp); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	if comp.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	var in landingComponentInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if in.Name != "" {
		comp.Name = in.Name
	}
	if len(in.Tree) > 0 {
		comp.Tree = datatypes.JSON(in.Tree)
	}
	if _, err := facades.Orm().Query().Where("id = ?", id).Update(&comp); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": comp})
}

// Destroy removes a component master. Existing instances in pages become
// unresolvable (renderer falls back to an empty placeholder).
func (c *LandingComponentController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var comp models.LandingComponent
	if err := facades.Orm().Query().Where("id = ?", id).First(&comp); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	if comp.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "component not found"})
	}
	if _, err := facades.Orm().Query().Where("id = ?", id).Delete(&models.LandingComponent{}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": true})
}
