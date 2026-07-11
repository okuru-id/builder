package routes

import (
	"encoding/json"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
	"okuru/app/models"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		var mode models.Setting
		if err := facades.Orm().Query().Where("key = ?", "landing_mode").First(&mode); err == nil && mode.Value == "custom" {
			var html models.Setting
			if err := facades.Orm().Query().Where("key = ?", "landing_template_html").First(&html); err == nil && html.Value != "" {
				content := html.Value
				if strings.Contains(content, "{{DATA}}") {
					var sections []models.LandingSection
					facades.Orm().Query().Where("is_active", true).Order("sort_order asc").Get(&sections)
					sectionMap := make(map[string]any, len(sections))
					for _, s := range sections {
						sectionMap[s.Type] = s.Content
					}
					data, _ := json.Marshal(sectionMap)
					content = strings.ReplaceAll(content, "{{DATA}}", string(data))
				}
				return ctx.Response().Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			}
		}
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
