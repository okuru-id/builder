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

	// Public API (no auth) — consumed by the React frontend.
	facades.Route().Prefix("api/v1").Group(func(router route.Router) {
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
			landingPageController := admin.NewLandingPageController()
			landingComponentController := admin.NewLandingComponentController()

			r.Get("/landing-pages", landingPageController.Index)
			r.Post("/landing-pages", landingPageController.Store)
			r.Get("/landing-pages/{id}", landingPageController.Show)
			r.Put("/landing-pages/{id}", landingPageController.Update)
			r.Delete("/landing-pages/{id}", landingPageController.Destroy)
			r.Post("/landing-pages/{id}/publish", landingPageController.Publish)
			r.Post("/landing-pages/{id}/duplicate", landingPageController.Duplicate)
			r.Get("/landing-pages/{id}/revisions", landingPageController.Revisions)
			r.Post("/landing-pages/{id}/revisions/{rid}/restore", landingPageController.RestoreRevision)

			r.Get("/landing-components", landingComponentController.Index)
			r.Post("/landing-components", landingComponentController.Store)
			r.Get("/landing-components/{id}", landingComponentController.Show)
			r.Put("/landing-components/{id}", landingComponentController.Update)
			r.Delete("/landing-components/{id}", landingComponentController.Destroy)

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
