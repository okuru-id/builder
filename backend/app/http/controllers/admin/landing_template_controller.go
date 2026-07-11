package admin

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type LandingTemplateController struct{}

func NewLandingTemplateController() *LandingTemplateController {
	return &LandingTemplateController{}
}

// Index lists all templates.
func (c *LandingTemplateController) Index(ctx http.Context) http.Response {
	var templates []models.LandingTemplate
	facades.Orm().Query().Order("id asc").Get(&templates)
	return ctx.Response().Success().Json(http.Json{"data": templates})
}

// Show returns a single template.
func (c *LandingTemplateController) Show(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}
	var tmpl models.LandingTemplate
	if err := facades.Orm().Query().Find(&tmpl, id); err != nil || tmpl.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "template not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": tmpl})
}

// Store creates a template.
func (c *LandingTemplateController) Store(ctx http.Context) http.Response {
	var input struct {
		Name        string               `json:"name"`
		Description string               `json:"description"`
		Preview     string               `json:"preview"`
		Sections    models.LandingTemplateItems `json:"sections"`
		HTML        string               `json:"html"`
	}
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	if input.Name == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "name is required"})
	}

	tmpl := models.LandingTemplate{
		Name:        input.Name,
		Description: input.Description,
		Preview:     input.Preview,
		Sections:    input.Sections,
		HTML:        input.HTML,
	}
	if err := facades.Orm().Query().Create(&tmpl); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": tmpl})
}

// Update updates a template.
func (c *LandingTemplateController) Update(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	var input struct {
		Name        string                     `json:"name"`
		Description string                     `json:"description"`
		Preview     string                     `json:"preview"`
		Sections    models.LandingTemplateItems `json:"sections"`
		HTML        string                     `json:"html"`
	}
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	var tmpl models.LandingTemplate
	if err := facades.Orm().Query().Find(&tmpl, id); err != nil || tmpl.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "template not found"})
	}

	tmpl.Name = input.Name
	tmpl.Description = input.Description
	tmpl.Preview = input.Preview
	tmpl.Sections = input.Sections
	tmpl.HTML = input.HTML

	if err := facades.Orm().Query().Save(&tmpl); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	return ctx.Response().Success().Json(http.Json{"data": tmpl})
}

// Destroy deletes a template.
func (c *LandingTemplateController) Destroy(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}
	facades.Orm().Query().Delete(&models.LandingTemplate{}, id)
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}

// Apply replaces all landing sections with the template's sections.
func (c *LandingTemplateController) Apply(ctx http.Context) http.Response {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid id"})
	}

	var tmpl models.LandingTemplate
	if err := facades.Orm().Query().Find(&tmpl, id); err != nil || tmpl.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "template not found"})
	}

	// Delete existing sections
	facades.Orm().Query().Where("id > ?", 0).Delete(&models.LandingSection{})

	// Create sections from template
	for i, raw := range tmpl.Sections {
		section, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		typeVal, _ := section["type"].(string)
		if typeVal == "" {
			continue
		}
		s := models.LandingSection{
			Type:      typeVal,
			SortOrder: i,
			IsActive:  true,
		}
		if content, ok := section["content"]; ok {
			if cmap, ok := content.(map[string]any); ok {
				s.Content = models.LandingContent(cmap)
			}
		}
		if isActive, ok := section["is_active"]; ok {
			if active, ok := isActive.(bool); ok {
				s.IsActive = active
			}
		}
		if err := facades.Orm().Query().Create(&s); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
		}
	}

	return ctx.Response().Success().Json(http.Json{"message": "Template applied"})
}
