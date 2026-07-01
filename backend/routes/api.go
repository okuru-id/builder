package routes

import (
	"github.com/goravel/framework/contracts/route"

	"okuru/app/facades"
	"okuru/app/http/controllers"
	"okuru/app/http/middleware"
)

func Api() {
	authController := controllers.NewAuthController()

	facades.Route().Prefix("admin/api").Group(func(router route.Router) {
		router.Post("/auth/login", authController.Login)
		router.Post("/auth/totp", authController.VerifyTotp)

		router.Middleware(middleware.JwtAuth()).Group(func(r route.Router) {
			r.Post("/auth/totp/setup", authController.SetupTotp)
			r.Post("/auth/totp/verify-setup", authController.VerifyTotpSetup)
			r.Get("/auth/me", authController.Me)
		})

		router.Middleware(middleware.JwtAuth(), middleware.TotpVerified()).Group(func(r route.Router) {
		})
	})
}
