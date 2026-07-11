package routes

import (
	"github.com/goravel/framework/contracts/route"

	"okuru/app/facades"
	"okuru/app/http/controllers"
	"okuru/app/http/controllers/admin"
	"okuru/app/http/middleware"
)

func Api() {
	authController := controllers.NewAuthController()
	contentController := controllers.NewContentController()
	landingController := controllers.NewLandingController()

	// Public API (no auth) — consumed by the React frontend.
	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
		router.Get("/landing", landingController.Index)
		router.Get("/content", contentController.Content)
		router.Get("/posts", contentController.Posts)
		router.Get("/posts/{slug}", contentController.Post)
		router.Get("/projects", contentController.Projects)
		router.Get("/open-source", contentController.OpenSource)
		router.Get("/categories", contentController.Categories)
		router.Get("/settings/public", contentController.PublicSettings)
		router.Post("/contact", contentController.Contact)
	})

	// Admin API.
	facades.Route().Prefix("admin/api").Group(func(router route.Router) {
		router.Post("/auth/login", authController.Login)
		router.Post("/auth/totp", authController.VerifyTotp)

		router.Middleware(middleware.JwtAuth()).Group(func(r route.Router) {
			r.Post("/auth/totp/setup", authController.SetupTotp)
			r.Post("/auth/totp/verify-setup", authController.VerifyTotpSetup)
			r.Get("/auth/me", authController.Me)
		})

		router.Middleware(middleware.JwtAuth(), middleware.TotpVerified()).Group(func(r route.Router) {
			postController := admin.NewPostController()
			projectController := admin.NewProjectController()
			openSourceController := admin.NewOpenSourceController()
			productController := admin.NewProductController()
			categoryController := admin.NewCategoryController()
			messageController := admin.NewMessageController()
			settingController := admin.NewSettingController()
			uploadController := admin.NewUploadController()
			landingTemplateController := admin.NewLandingTemplateController()
			landingSectionController := admin.NewLandingSectionController()

			r.Get("/landing-templates", landingTemplateController.Index)
			r.Post("/landing-templates", landingTemplateController.Store)
			r.Get("/landing-templates/{id}", landingTemplateController.Show)
			r.Put("/landing-templates/{id}", landingTemplateController.Update)
			r.Delete("/landing-templates/{id}", landingTemplateController.Destroy)
			r.Post("/landing-templates/{id}/apply", landingTemplateController.Apply)

			r.Get("/landing-sections", landingSectionController.Index)
			r.Post("/landing-sections", landingSectionController.Store)
			r.Put("/landing-sections/{id}", landingSectionController.Update)
			r.Delete("/landing-sections/{id}", landingSectionController.Destroy)
			r.Patch("/landing-sections/{id}/toggle", landingSectionController.Toggle)
			r.Patch("/landing-sections/{id}/sort", landingSectionController.Sort)

			r.Get("/posts", postController.Index)
			r.Post("/posts", postController.Store)
			r.Get("/posts/{id}", postController.Show)
			r.Put("/posts/{id}", postController.Update)
			r.Delete("/posts/{id}", postController.Destroy)

			r.Get("/projects", projectController.Index)
			r.Post("/projects", projectController.Store)
			r.Get("/projects/{id}", projectController.Show)
			r.Put("/projects/{id}", projectController.Update)
			r.Delete("/projects/{id}", projectController.Destroy)

			r.Get("/open-source", openSourceController.Index)
			r.Post("/open-source", openSourceController.Store)
			r.Get("/open-source/{id}", openSourceController.Show)
			r.Put("/open-source/{id}", openSourceController.Update)
			r.Delete("/open-source/{id}", openSourceController.Destroy)

			r.Get("/products", productController.Index)
			r.Post("/products", productController.Store)
			r.Get("/products/{id}", productController.Show)
			r.Put("/products/{id}", productController.Update)
			r.Delete("/products/{id}", productController.Destroy)

			r.Get("/categories", categoryController.Index)
			r.Post("/categories", categoryController.Store)
			r.Get("/categories/{id}", categoryController.Show)
			r.Put("/categories/{id}", categoryController.Update)
			r.Delete("/categories/{id}", categoryController.Destroy)

			r.Get("/messages", messageController.Index)
			r.Get("/messages/{id}", messageController.Show)
			r.Post("/messages/{id}/read", messageController.MarkRead)
			r.Post("/messages/{id}/archive", messageController.Archive)
			r.Delete("/messages/{id}", messageController.Destroy)

			r.Get("/settings", settingController.Index)
			r.Put("/settings", settingController.Update)
			r.Delete("/settings/{key}", settingController.Destroy)

			r.Post("/upload", uploadController.Store)
		})
	})
}
