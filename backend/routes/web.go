package routes

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().File("./public/index.html")
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
