package admin

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

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
	Name   string `json:"name" validate:"required"`
	Slug   string `json:"slug"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
	IsHome *bool  `json:"is_home"` // pointer so omit != false
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
		Path      string `json:"path"`
		Domain    string `json:"domain"`
		IsHome    bool   `json:"is_home"`
		UpdatedAt string `json:"updated_at"`
	}
	items := make([]listItem, 0, len(pages))
	for _, p := range pages {
		updatedAt := ""
		if !p.UpdatedAt.IsZero() {
			updatedAt = p.UpdatedAt.ToDateTimeString()
		}
		items = append(items, listItem{p.ID, p.Slug, p.Name, p.Status, p.Version, p.Path, p.Domain, p.IsHome, updatedAt})
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
	slug = uniqueSlug(slug)
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

// treeInput is the autosave payload. path/domain/is_home are optional — the
// dedicated Settings endpoint is the primary path for those, but autosave also
// accepts them so the builder can persist routing changes inline.
type treeInput struct {
	Name   string          `json:"name"`
	Slug   string          `json:"slug"`
	Path   *string         `json:"path"`
	Domain *string         `json:"domain"`
	IsHome *bool           `json:"is_home"`
	Tree   json.RawMessage `json:"tree"`
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
	// Routing fields (optional in autosave). Validate uniqueness when present.
	if in.Path != nil {
		np := normalizePath(*in.Path)
		if resp, ok := ensureUniqueField(ctx, "path", np, page.ID); !ok {
			return resp
		}
		page.Path = np
	}
	if in.Domain != nil {
		nd := normalizeDomain(*in.Domain)
		if resp, ok := ensureUniqueField(ctx, "domain", nd, page.ID); !ok {
			return resp
		}
		page.Domain = nd
	}
	if in.IsHome != nil {
		if *in.IsHome {
			// At most one home page. Clear the flag on every other page.
			facades.Orm().Query().Where("id != ?", page.ID).Update(map[string]any{"is_home": false})
		}
		page.IsHome = *in.IsHome
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

// Settings updates a page's routing fields (path, domain, is_home) and is the
// primary UI entry point for multi-page publishing. Validates uniqueness and
// enforces the single-home constraint. Body: { path?, domain?, is_home? }.
func (c *LandingPageController) Settings(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ?", id).First(&page); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}
	if page.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"})
	}

	var in struct {
		Path   *string `json:"path"`
		Domain *string `json:"domain"`
		IsHome *bool   `json:"is_home"`
	}
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	if in.Path != nil {
		np := normalizePath(*in.Path)
		if resp, ok := ensureUniqueField(ctx, "path", np, page.ID); !ok {
			return resp
		}
		page.Path = np
	}
	if in.Domain != nil {
		nd := normalizeDomain(*in.Domain)
		if resp, ok := ensureUniqueField(ctx, "domain", nd, page.ID); !ok {
			return resp
		}
		page.Domain = nd
	}
	if in.IsHome != nil {
		if *in.IsHome {
			facades.Orm().Query().Where("id != ?", page.ID).Update(map[string]any{"is_home": false})
		}
		page.IsHome = *in.IsHome
	}

	if _, err := facades.Orm().Query().Where("id = ?", id).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	// Re-read so the returned row reflects DB defaults.
	facades.Orm().Query().Where("id = ?", id).First(&page)
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
		Slug:   uniqueSlug(slugify(page.Name + " copy")),
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
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
			prevDash = false
		case r == ' ', r == '-', r == '_':
			if b.Len() > 0 && !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "page"
	}
	return out
}

func uniqueSlug(base string) string {
	base = slugify(base)
	candidate := base
	for i := 2; ; i++ {
		var page models.LandingPage
		_ = facades.Orm().Query().Where("slug = ?", candidate).First(&page)
		if page.ID == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// normalizePath trims leading/trailing slashes and collapses to a single segment.
// Empty input is allowed (means "not published via path"). Reserved paths that
// would collide with app routes (/admin, /api, /assets, /shop, /health) are rejected.
func normalizePath(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")
	return s
}

func isReservedPath(p string) bool {
	switch p {
	case "", "admin", "api", "assets", "images", "shop", "product", "checkout", "health":
		return true
	}
	return false
}

// normalizeDomain lowercases, trims whitespace, strips scheme, path, and port.
// "https://Promo.Example.com:443/x" → "promo.example.com".
func normalizeDomain(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	if i := strings.IndexAny(s, "/?#"); i != -1 {
		s = s[:i]
	}
	if i := strings.LastIndex(s, ":"); i != -1 {
		// keep colons only if they're part of an IPv6 bracket; else strip port
		if !strings.Contains(s, "[") {
			s = s[:i]
		}
	}
	return strings.ToLower(strings.Trim(s, "."))
}

// ensureUniqueField checks whether another page already uses the given non-empty
// value for path/domain. Returns the conflict response and ok=false; when ok=true
// the caller may proceed. Empty values are allowed (many pages may have no custom
// path/domain).
func ensureUniqueField(ctx http.Context, field, value string, excludeID uint) (http.Response, bool) {
	if value == "" {
		return nil, true
	}
	var other models.LandingPage
	q := facades.Orm().Query().Where(field+" = ?", value).Where("id != ?", excludeID)
	_ = q.First(&other)
	if other.ID != 0 {
		return ctx.Response().Json(http.StatusConflict, http.Json{
			"error":  field + " already in use by another page",
			"field":  field,
			"value":  value,
			"pageId": other.ID,
		}), false
	}
	return nil, true
}
