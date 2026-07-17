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
	"okuru/resources"
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
	facades.Route().Static("assets", "./public/assets")
	facades.Route().Static("images", "./public/images")
	facades.Route().Static("public", "./public")

	facades.Route().Fallback(func(ctx http.Context) http.Response {
		path := strings.TrimPrefix(ctx.Request().Path(), "/")

		if strings.HasPrefix(path, "assets/") || strings.HasPrefix(path, "images/") || strings.HasPrefix(path, "public/") {
			if _, err := os.Stat("./public/" + path); err == nil {
				return ctx.Response().File("./public/" + path)
			}
			return ctx.Response().String(http.StatusNotFound, "Not Found")
		}

		// Root-level files (favicon, icons, etc.)
		if !strings.Contains(path, "/") {
			if _, err := os.Stat("./public/" + path); err == nil {
				return ctx.Response().File("./public/" + path)
			}
		}

		if strings.HasPrefix(path, "admin") {
			return ctx.Response().File("./public/admin.html")
		}

		html, status := resolveStorefront(ctx, path)
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(status, html)
	})

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

	// 2b. Unknown host: a domain pointed here in DNS but not bound to any
	// page. Serve a welcome page instead of leaking the okuru.id marketing
	// landing. Primary app domains (okuru.id, www, subdomains, localhost)
	// fall through to the normal home/default flow.
	// ponytail: extend via APP_HOSTS env (comma-separated) when adding subs.
	if host != "" && !isPrimaryHost(host) {
		return welcomePageHTML(host), http.StatusOK
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

// isPrimaryHost reports whether host belongs to the app itself (marketing,
// api, shop, builder) or is a parked/unknown domain. Override the list via
// the APP_HOSTS env var (comma-separated) without touching code.
// ponytail: env over config file — fewer moving parts, redeploy = update.
var primaryHosts = []string{
	"okuru.id", "www.okuru.id", "shop.okuru.id", "api.okuru.id",
	"builder.okuru.id", "localhost", "127.0.0.1",
}

func isPrimaryHost(host string) bool {
	if extra := os.Getenv("APP_HOSTS"); extra != "" {
		for _, h := range strings.Split(extra, ",") {
			if strings.TrimSpace(strings.ToLower(h)) == host {
				return true
			}
		}
	}
	for _, h := range primaryHosts {
		if h == host {
			return true
		}
	}
	return false
}

// welcomePageHTML is the placeholder shown when a domain resolves to this
// server but no published page is bound to it yet. Sourced from
// resources/welcome.html via go:embed — edit the file directly.
// ponytail: embed over runtime read — no IO, no missing-file path to handle.
func welcomePageHTML(_ string) string {
	return resources.WelcomePage
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
