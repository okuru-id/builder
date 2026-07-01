package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type ContentController struct{}

func NewContentController() *ContentController {
	return &ContentController{}
}

// GET /api/v1/content — batch endpoint returning all site content.
func (c *ContentController) Content(ctx http.Context) http.Response {
	settings := c.loadSettings()

	var projects []models.Project
	facades.Orm().Query().OrderBy("sort_order", "asc").Get(&projects)

	var openSource []models.OpenSourceProject
	facades.Orm().Query().OrderBy("sort_order", "asc").Get(&openSource)

	var categories []models.Category
	facades.Orm().Query().OrderBy("id", "asc").Get(&categories)

	var posts []models.Post
	facades.Orm().Query().
		Where("status = ?", "published").
		OrderBy("published_at", "desc").
		Limit(5).
		Get(&posts)

	return ctx.Response().Success().Json(http.Json{
		"settings":    settings,
		"projects":    projects,
		"open_source": openSource,
		"categories":  categories,
		"posts":       posts,
	})
}

// GET /api/v1/posts — list published posts with optional ?category and ?page.
func (c *ContentController) Posts(ctx http.Context) http.Response {
	page, _ := strconv.Atoi(ctx.Request().Query("page", "1"))
	if page < 1 {
		page = 1
	}
	const perPage = 10

	query := facades.Orm().Query().
		Where("status = ?", "published").
		OrderBy("published_at", "desc")

	if category := ctx.Request().Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	var posts []models.Post
	query.Offset((page - 1) * perPage).Limit(perPage).Get(&posts)

	return ctx.Response().Success().Json(http.Json{
		"data": posts,
		"page": page,
	})
}

// GET /api/v1/posts/{slug} — single published post by slug.
func (c *ContentController) Post(ctx http.Context) http.Response {
	slug := ctx.Request().Route("slug")

	var post models.Post
	if err := facades.Orm().Query().
		Where("slug = ?", slug).
		Where("status = ?", "published").
		First(&post); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Post not found"})
	}

	return ctx.Response().Success().Json(http.Json{"data": post})
}

// GET /api/v1/projects — list all projects ordered by sort_order.
func (c *ContentController) Projects(ctx http.Context) http.Response {
	var projects []models.Project
	facades.Orm().Query().OrderBy("sort_order", "asc").Get(&projects)

	return ctx.Response().Success().Json(http.Json{"data": projects})
}

// GET /api/v1/open-source — list all OS projects ordered by sort_order.
func (c *ContentController) OpenSource(ctx http.Context) http.Response {
	var projects []models.OpenSourceProject
	facades.Orm().Query().OrderBy("sort_order", "asc").Get(&projects)

	return ctx.Response().Success().Json(http.Json{"data": projects})
}

// GET /api/v1/categories — list all categories.
func (c *ContentController) Categories(ctx http.Context) http.Response {
	var categories []models.Category
	facades.Orm().Query().OrderBy("id", "asc").Get(&categories)

	return ctx.Response().Success().Json(http.Json{"data": categories})
}

// GET /api/v1/settings/public — public settings as key-value map.
func (c *ContentController) PublicSettings(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{"data": c.loadSettings()})
}

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// POST /api/v1/contact — submit contact form.
func (c *ContentController) Contact(ctx http.Context) http.Response {
	var input contactRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "invalid request body"})
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Message = strings.TrimSpace(input.Message)
	if input.Name == "" || input.Email == "" || input.Message == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "name, email, and message are required"})
	}

	message := models.Message{
		Name:    input.Name,
		Email:   input.Email,
		Content: input.Message,
		Status:  "unread",
	}
	if err := facades.Orm().Query().Create(&message); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": "failed to save message"})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": message})
}

// loadSettings fetches all settings as a key-value map.
func (c *ContentController) loadSettings() map[string]string {
	var settings []models.Setting
	facades.Orm().Query().Get(&settings)

	out := make(map[string]string, len(settings))
	for _, s := range settings {
		out[s.Key] = s.Value
	}
	return out
}
