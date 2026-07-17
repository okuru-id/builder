package routes

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
	"okuru/app/models"
	"okuru/app/services"
)

func Web() {
	// Storefront root + preview. Path-based pages (/promo) are handled by Fallback.
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		html, status := resolveStorefront(ctx, "")
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(status, html)
	})

	// Admin SPA: serve the shell for /admin and any /admin/* client route.
	facades.Route().Get("/admin/", func(ctx http.Context) http.Response {
		return ctx.Response().File("./public/admin.html")
	})
	// SPA fallback + storefront path dispatch:
	//   - /admin/*           → admin shell (client router)
	//   - any other segment   → resolve as a published page by path
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		path := strings.TrimPrefix(ctx.Request().Path(), "/")
		if strings.HasPrefix(path, "admin") {
			return ctx.Response().File("./public/admin.html")
		}
		html, status := resolveStorefront(ctx, path)
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(status, html)
	})

	facades.Route().Static("assets", "./public/assets")
	facades.Route().Static("images", "./public/images")
	facades.Route().Static("public", "./public")

	userController := controllers.NewUserController()
	facades.Route().Get("/users", userController.Index)

	shopController := controllers.NewShopController()
	facades.Route().Get("/shop", shopController.Index)
	facades.Route().Get("/product/{slug}", shopController.Product)
	facades.Route().Post("/checkout/{slug}", shopController.Checkout)

	facades.Route().Get("/health", func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusOK, http.Json{
			"status": "ok",
		})
	})
}

// resolveStorefront picks which published HTML to serve for a request.
//
// Resolution priority:
//  1. ?preview=<id>  → render the draft tree of that page (no publish required).
//  2. Host header     → page whose Domain matches (published only).
//  3. path == ""      → page flagged is_home (published), else the app's own
//     marketing landing (defaultStorefrontHTML).
//  4. path != ""      → page whose Path equals the segment (published).
//     Unknown path → 404.
//
// Returns (body, status). Callers wrap with Content-Type.
func resolveStorefront(ctx http.Context, path string) (string, int) {
	// 1. Preview mode (draft render). Allowed on any path/host.
	if previewID := ctx.Request().Query("preview", ""); previewID != "" {
		var page models.LandingPage
		if err := facades.Orm().Query().Where("id = ?", previewID).First(&page); err != nil || page.ID == 0 {
			return "Halaman tidak ditemukan", http.StatusNotFound
		}
		return renderPreviewHTML(page, ctx.Request().Host()), http.StatusOK
	}

	host := normalizeHost(ctx.Request().Host())

	// 2. Custom-domain match.
	if host != "" {
		var page models.LandingPage
		_ = facades.Orm().Query().
			Where("domain = ?", host).
			Where("status = ?", "published").
			First(&page)
		if page.ID != 0 && page.PublishedHTML != "" {
			return page.PublishedHTML, http.StatusOK
		}
	}

	// 3. Home.
	if path == "" {
		var page models.LandingPage
		_ = facades.Orm().Query().
			Where("is_home = ?", true).
			Where("status = ?", "published").
			First(&page)
		if page.ID != 0 && page.PublishedHTML != "" {
			return page.PublishedHTML, http.StatusOK
		}
		return defaultStorefrontHTML(ctx.Request().Host()), http.StatusOK
	}

	// 4. Path match.
	var page models.LandingPage
	_ = facades.Orm().Query().
		Where("path = ?", path).
		Where("status = ?", "published").
		First(&page)
	if page.ID != 0 && page.PublishedHTML != "" {
		return page.PublishedHTML, http.StatusOK
	}
	return "404 — halaman tidak ditemukan", http.StatusNotFound
}

// normalizeHost lowercases the Host header and strips the port.
func normalizeHost(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	if i := strings.LastIndex(h, ":"); i != -1 && !strings.Contains(h, "[") {
		h = h[:i]
	}
	return strings.Trim(h, ".")
}

// renderPreviewHTML converts a page's working tree into a full HTML document
// using the codegen pipeline. Used for ?preview= mode.
func renderPreviewHTML(page models.LandingPage, host string) string {
	var m map[string]any
	if err := json.Unmarshal(page.Tree, &m); err != nil {
		return defaultStorefrontHTML(host)
	}
	resolved := services.ResolveComponentInstances(m)
	return services.NewLandingCodegen().Generate(resolved, page.Name+" (preview)")
}
// defaultStorefrontHTML serves the built Vue marketing landing (dist/index.html).
// In dev, Vite serves `/` directly with HMR — this path is only hit in prod.
// ponytail: reuse built assets instead of duplicating markup in Go.
func defaultStorefrontHTML(host string) string {
	if b, err := os.ReadFile("./public/index.html"); err == nil {
		return string(b)
	}
	return "<!doctype html><html><body>okuru.id</body></html>"
}
