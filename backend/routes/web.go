package routes

import (
	"encoding/json"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
	"okuru/app/models"
	"okuru/app/services"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		// Preview mode: render the draft (working) tree for a specific page.
		previewID := ctx.Request().Query("preview", "")
		if previewID != "" {
			var page models.LandingPage
			if err := facades.Orm().Query().Where("id = ?", previewID).First(&page); err != nil || page.ID == 0 {
				return ctx.Response().String(http.StatusNotFound, "Halaman tidak ditemukan")
			}
			// Render the working tree (not published) so the preview shows current edits.
			html := renderPreviewHTML(page, ctx.Request().Host())
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(http.StatusOK, html)
		}

		// Storefront: serve the most recently published landing page's cached HTML.
		// Falls back to a placeholder when nothing is published yet.
		var page models.LandingPage
		if err := facades.Orm().Query().Where("status = ?", "published").Order("updated_at desc").First(&page); err != nil || page.ID == 0 {
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(http.StatusOK, defaultStorefrontHTML(ctx.Request().Host()))
		}
		if page.PublishedHTML != "" {
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(http.StatusOK, page.PublishedHTML)
		}
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(http.StatusOK, defaultStorefrontHTML(ctx.Request().Host()))
	})

	// Admin SPA: serve the shell for /admin and any /admin/* client route.
	facades.Route().Get("/admin/", func(ctx http.Context) http.Response {
		return ctx.Response().File("./public/admin.html")
	})
	// SPA fallback: any unmatched /admin/* client route serves the admin shell.
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		if strings.HasPrefix(ctx.Request().Path(), "/admin") {
			return ctx.Response().File("./public/admin.html")
		}
		return ctx.Response().String(http.StatusNotFound, "Not Found")
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

// defaultStorefrontHTML is the placeholder shown when no landing page is published.
// ponytail: static string, no template engine. Replace with a real template only
// when the storefront needs server-side data injection (Phase 7+).
func defaultStorefrontHTML(host string) string {
	// Strip port for display: "okuru.id:443" -> "okuru.id".
	name := host
	if i := strings.LastIndex(host, ":"); i != -1 {
		name = host[:i]
	}
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>` + name + `</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="min-h-screen flex items-center justify-center bg-neutral-50 text-neutral-500">
  <div class="text-center">
    <h1 class="text-2xl font-semibold text-neutral-900">` + name + `</h1>
    <p class="mt-2">Landing page belum dipublikasi.</p>
  </div>
</body>
</html>`
}
