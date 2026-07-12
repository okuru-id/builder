package routes

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
	"okuru/app/models"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		// Storefront: serve the most recently published landing page's cached HTML.
		// Falls back to a placeholder when nothing is published yet.
		var page models.LandingPage
		if err := facades.Orm().Query().Where("status = ?", "published").Order("updated_at desc").First(&page); err != nil || page.ID == 0 {
			return ctx.Response().String(http.StatusOK, defaultStorefrontHTML())
		}
		if page.PublishedHTML != "" {
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(http.StatusOK, page.PublishedHTML)
		}
		return ctx.Response().String(http.StatusOK, defaultStorefrontHTML())
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

// defaultStorefrontHTML is the placeholder shown when no landing page is published.
// ponytail: static string, no template engine. Replace with a real template only
// when the storefront needs server-side data injection (Phase 7+).
func defaultStorefrontHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>okuru.id</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="min-h-screen flex items-center justify-center bg-neutral-50 text-neutral-500">
  <div class="text-center">
    <h1 class="text-2xl font-semibold text-neutral-900">okuru.id</h1>
    <p class="mt-2">Landing page belum dipublikasi.</p>
  </div>
</body>
</html>`
}
