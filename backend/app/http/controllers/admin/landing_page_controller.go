package admin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"gorm.io/datatypes"

	"okuru/app/facades"
	"okuru/app/models"
	"okuru/app/services"
)

type LandingPageController struct{}

func NewLandingPageController() *LandingPageController {
	return &LandingPageController{}
}

type landingPageInput struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug"`
}

// defaultRootTree returns the initial empty tree for a new page.
func defaultRootTree() datatypes.JSON {
	tree := map[string]any{
		"root": map[string]any{
			"id":       "root",
			"type":     "frame",
			"name":     "Page",
			"props":    map[string]any{},
			"classes":  []string{"min-h-screen", "bg-white", "text-neutral-900"},
			"children": []any{},
		},
	}
	b, _ := json.Marshal(tree)
	return b
}

// Index lists all landing pages (without heavy tree payload).
func (c *LandingPageController) Index(ctx http.Context) http.Response {
	var pages []models.LandingPage
	if err := facades.Orm().Query().Order("id desc").Get(&pages); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	type listItem struct {
		ID        uint   `json:"id"`
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Status    string `json:"status"`
		Version   int    `json:"version"`
		UpdatedAt string `json:"updated_at"`
	}
	items := make([]listItem, 0, len(pages))
	for _, p := range pages {
		updatedAt := ""
		if !p.UpdatedAt.IsZero() {
			updatedAt = p.UpdatedAt.ToDateTimeString()
		}
		items = append(items, listItem{p.ID, p.Slug, p.Name, p.Status, p.Version, updatedAt})
	}
	return ctx.Response().Success().Json(http.Json{"data": items})
}

// Show returns a single page with its full tree.
func (c *LandingPageController) Show(ctx http.Context) http.Response {
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", ctx.Request().Input("id")).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// Store creates a new draft page.
func (c *LandingPageController) Store(ctx http.Context) http.Response {
	var in landingPageInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	slug := in.Slug
	if slug == "" {
		slug = slugify(in.Name)
	}
	page := models.LandingPage{
		Name:   in.Name,
		Slug:   slug,
		Status: "draft",
		Tree:   defaultRootTree(),
	}
	if err := facades.Orm().Query().Create(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// treeInput is the autosave payload.
type treeInput struct {
	Name string          `json:"name"`
	Slug string          `json:"slug"`
	Tree json.RawMessage `json:"tree"`
}

// Update autosaves the working tree, bumps version, and logs a revision.
func (c *LandingPageController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", id).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}

	var in treeInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	// Save revision of the current tree before overwriting.
	rev := models.LandingPageRevision{
		LandingPageID: page.ID,
		Tree:          page.Tree,
		Message:       fmt.Sprintf("autosave v%d", page.Version),
	}
	if err := facades.Orm().Query().Create(&rev); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	page.Tree = datatypes.JSON(in.Tree)
	if in.Name != "" {
		page.Name = in.Name
	}
	if in.Slug != "" {
		page.Slug = in.Slug
	}
	if _, err := facades.Orm().Query().Where("id = ?", id).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page, "revision": rev.ID})
}

// Publish snapshots the working tree into published_tree and marks status=published.
func (c *LandingPageController) Publish(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", id).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	page.PublishedTree = page.Tree
	page.Status = "published"
	page.Version = page.Version + 1

	// Run tree → HTML codegen synchronously at publish time. Cached in published_html,
	// served raw by the storefront so the published page never depends on the editor.
	// Component instances are resolved to their master trees first, so published HTML
	// is self-contained even if a master is later edited or deleted.
	resolved := services.ResolveComponentInstances(treeToMap(page.Tree))
	if html, err := renderHTMLFromMap(resolved, page.Name); err == nil {
		page.PublishedHTML = html
	}

	if _, err := facades.Orm().Query().Where("id = ?", id).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// renderHTML converts a datatypes.JSON tree into a full HTML document via the
// LandingCodegen service.
func renderHTML(tree datatypes.JSON, title string) (string, error) {
	m := treeToMap(tree)
	return renderHTMLFromMap(m, title)
}

func renderHTMLFromMap(m map[string]any, title string) (string, error) {
	return services.NewLandingCodegen().Generate(m, title), nil
}

func treeToMap(tree datatypes.JSON) map[string]any {
	var m map[string]any
	if err := json.Unmarshal(tree, &m); err != nil {
		return map[string]any{"root": map[string]any{}}
	}
	return m
}

// Revisions lists the page revision history (newest first).
func (c *LandingPageController) Revisions(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var revs []models.LandingPageRevision
	if err := facades.Orm().Query().Where("landing_page_id = ?", id).Order("id desc").Limit(100).Get(&revs); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": revs})
}

// RestoreRevision rolls the working tree back to a prior revision and bumps version.
func (c *LandingPageController) RestoreRevision(ctx http.Context) http.Response {
	pageID := ctx.Request().Input("id")
	revID := ctx.Request().Input("rid")

	var rev models.LandingPageRevision
	if err := facades.Orm().Query().Where("id = ? AND landing_page_id = ?", revID, pageID).First(&rev); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "revision not found"})
	}
	if rev.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "revision not found"})
	}

	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", pageID).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	page.Tree = rev.Tree
	page.Version = page.Version + 1
	if _, err := facades.Orm().Query().Where("id = ?", pageID).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// Destroy deletes a landing page and its revisions.
func (c *LandingPageController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", id).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	// Delete associated revisions first.
	facades.Orm().Query().Where("landing_page_id = ?", id).Delete(&models.LandingPageRevision{})
	facades.Orm().Query().Where("id = ?", id).Delete(&models.LandingPage{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}

// Duplicate creates a copy of an existing page with "(copy)" suffix.
func (c *LandingPageController) Duplicate(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", id).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	copy := models.LandingPage{
		Name:   page.Name + " (copy)",
		Slug:   slugify(page.Name + " copy"),
		Status: "draft",
		Tree:   page.Tree,
	}
	if err := facades.Orm().Query().Create(&copy); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": copy})
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}
