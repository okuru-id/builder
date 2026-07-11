package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type LandingController struct{}

func NewLandingController() *LandingController {
	return &LandingController{}
}

// Index returns all active landing sections grouped by type.
func (c *LandingController) Index(ctx http.Context) http.Response {
	var sections []models.LandingSection
	if err := facades.Orm().Query().
		Where("is_active", true).
		Order("sort_order asc").
		Get(&sections); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	result := make(map[string]any, len(sections))
	for _, s := range sections {
		result[s.Type] = s.Content
	}

	return ctx.Response().Success().Json(http.Json{"data": result})
}
