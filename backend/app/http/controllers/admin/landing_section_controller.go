package admin

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type LandingSectionController struct{}

func NewLandingSectionController() *LandingSectionController {
	return &LandingSectionController{}
}

// Index lists all landing sections.
func (c *LandingSectionController) Index(ctx http.Context) http.Response {
	var sections []models.LandingSection
	facades.Orm().Query().Order("sort_order asc").Get(&sections)
	return ctx.Response().Success().Json(http.Json{"data": sections})
}

type storeLandingSectionRequest struct {
	Type    string                `json:"type"`
	Content models.LandingContent `json:"content"`
}

// Store creates a new landing section.
func (c *LandingSectionController) Store(ctx http.Context) http.Response {
	var input storeLandingSectionRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if input.Type == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "type is required"})
	}

	var maxSort models.LandingSection
	facades.Orm().Query().Order("sort_order desc").First(&maxSort)

	section := models.LandingSection{
		Type:      input.Type,
		Content:   input.Content,
		SortOrder: maxSort.SortOrder + 1,
		IsActive:  true,
	}
	if err := facades.Orm().Query().Create(&section); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	return ctx.Response().Success().Json(http.Json{"data": section})
}

// Update updates a landing section's content.
func (c *LandingSectionController) Update(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	var input storeLandingSectionRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	var section models.LandingSection
	if err := facades.Orm().Query().Find(&section, id); err != nil || section.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "section not found"})
	}

	section.Content = input.Content
	if err := facades.Orm().Query().Save(&section); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	return ctx.Response().Success().Json(http.Json{"data": section})
}

// Destroy deletes a landing section.
func (c *LandingSectionController) Destroy(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	facades.Orm().Query().Delete(&models.LandingSection{}, id)
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}

// Toggle toggles is_active on a landing section.
func (c *LandingSectionController) Toggle(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	var section models.LandingSection
	if err := facades.Orm().Query().Find(&section, id); err != nil || section.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "section not found"})
	}

	section.IsActive = !section.IsActive
	facades.Orm().Query().Save(&section)

	return ctx.Response().Success().Json(http.Json{"data": section})
}

// Sort updates the sort_order of a landing section.
func (c *LandingSectionController) Sort(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	var input struct {
		SortOrder int `json:"sort_order"`
	}
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	if _, err := facades.Orm().Query().Model(&models.LandingSection{}).Where("id", id).Update("sort_order", input.SortOrder); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	return ctx.Response().Success().Json(http.Json{"updated": true})
}
