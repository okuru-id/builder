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
	IsHome *bool  `json:"is_home"`
}

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

// ownPage fetches a page scoped to the caller. Returns ok=false with a 404
// response already built when missing/not-owned.
func ownPage(ctx http.Context, id string) (models.LandingPage, http.Response, bool) {
	var page models.LandingPage
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", id, currentUserID(ctx)).First(&page); err != nil || page.ID == 0 {
		return page, ctx.Response().Json(http.StatusNotFound, http.Json{"error": "page not found"}), false
	}
	return page, nil, true
}

// Index lists the caller's own pages (without heavy tree payload).
func (c *LandingPageController) Index(ctx http.Context) http.Response {
	var pages []models.LandingPage
	if err := facades.Orm().Query().Where("user_id = ?", currentUserID(ctx)).Order("id desc").Get(&pages); err != nil {
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
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// Store creates a new draft page owned by the caller.
func (c *LandingPageController) Store(ctx http.Context) http.Response {
	var in landingPageInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	slug := in.Slug
	if slug == "" {
		slug = slugify(in.Name)
	}
	slug = uniqueSlug(slug, currentUserID(ctx))
	page := models.LandingPage{
		Name:   in.Name,
		Slug:   slug,
		Status: "draft",
		Tree:   defaultRootTree(),
		UserID: currentUserID(ctx),
	}
	if err := facades.Orm().Query().Create(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

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
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}

	var in treeInput
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

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
	if in.Path != nil {
		np := normalizePath(*in.Path)
		if r, ok := ensureUniqueField(ctx, "path", np, page.ID); !ok {
			return r
		}
		page.Path = np
	}
	if in.Domain != nil {
		nd := normalizeDomain(*in.Domain)
		if r, ok := ensureUniqueField(ctx, "domain", nd, page.ID); !ok {
			return r
		}
		page.Domain = nd
	}
	if in.IsHome != nil {
		if *in.IsHome {
			// At most one home page globally (single-tenant storefront).
			facades.Orm().Query().Where("id != ?", page.ID).Update(map[string]any{"is_home": false})
		}
		page.IsHome = *in.IsHome
	}
	if _, err := facades.Orm().Query().Where("id = ?", page.ID).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page, "revision": rev.ID})
}

// Publish snapshots the working tree into published_tree and marks status=published.
func (c *LandingPageController) Publish(ctx http.Context) http.Response {
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	page.PublishedTree = page.Tree
	page.Status = "published"
	page.Version = page.Version + 1

	resolved := services.ResolveComponentInstances(treeToMap(page.Tree))
	if html, err := renderHTMLFromMap(resolved, page.Name); err == nil {
		page.PublishedHTML = html
	}

	if _, err := facades.Orm().Query().Where("id = ?", page.ID).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": page})
}

// Settings updates a page's routing fields (path, domain, is_home).
func (c *LandingPageController) Settings(ctx http.Context) http.Response {
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
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
		if r, ok := ensureUniqueField(ctx, "path", np, page.ID); !ok {
			return r
		}
		page.Path = np
	}
	if in.Domain != nil {
		nd := normalizeDomain(*in.Domain)
		if r, ok := ensureUniqueField(ctx, "domain", nd, page.ID); !ok {
			return r
		}
		page.Domain = nd
	}
	if in.IsHome != nil {
		if *in.IsHome {
			facades.Orm().Query().Where("id != ?", page.ID).Update(map[string]any{"is_home": false})
		}
		page.IsHome = *in.IsHome
	}

	if _, err := facades.Orm().Query().Where("id = ?", page.ID).Update(&page); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	facades.Orm().Query().Where("id = ?", page.ID).First(&page)
	return ctx.Response().Success().Json(http.Json{"data": page})
}

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

// Revisions lists the page revision history (newest first). Scoped to the
// caller's pages.
func (c *LandingPageController) Revisions(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	if _, resp, ok := ownPage(ctx, id); !ok {
		return resp
	}
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

	page, resp, ok := ownPage(ctx, pageID)
	if !ok {
		return resp
	}

	var rev models.LandingPageRevision
	if err := facades.Orm().Query().Where("id = ? AND landing_page_id = ?", revID, pageID).First(&rev); err != nil || rev.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "revision not found"})
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
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	facades.Orm().Query().Where("landing_page_id = ?", page.ID).Delete(&models.LandingPageRevision{})
	facades.Orm().Query().Where("id = ?", page.ID).Delete(&models.LandingPage{})
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}

// Duplicate creates a copy of an existing page with "(copy)" suffix.
func (c *LandingPageController) Duplicate(ctx http.Context) http.Response {
	page, resp, ok := ownPage(ctx, ctx.Request().Input("id"))
	if !ok {
		return resp
	}
	uid := currentUserID(ctx)
	copy := models.LandingPage{
		Name:   page.Name + " (copy)",
		Slug:   uniqueSlug(slugify(page.Name+" copy"), uid),
		Status: "draft",
		Tree:   page.Tree,
		UserID: uid,
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

// uniqueSlug makes a slug unique within the caller's own pages.
func uniqueSlug(base string, uid uint) string {
	base = slugify(base)
	candidate := base
	for i := 2; ; i++ {
		var page models.LandingPage
		_ = facades.Orm().Query().Where("slug = ? AND user_id = ?", candidate, uid).First(&page)
		if page.ID == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

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

func normalizeDomain(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	if i := strings.IndexAny(s, "/?#"); i != -1 {
		s = s[:i]
	}
	if i := strings.LastIndex(s, ":"); i != -1 {
		if !strings.Contains(s, "[") {
			s = s[:i]
		}
	}
	return strings.ToLower(strings.Trim(s, "."))
}

// ensureUniqueField checks GLOBAL uniqueness of path/domain — the public
// storefront resolves pages by these fields without a user context, so they
// must not collide across owners.
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
