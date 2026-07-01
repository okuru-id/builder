package middleware

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
)

func JwtAuth() http.Middleware {
	return func(ctx http.Context) {
		authHeader := ctx.Request().Header("Authorization", "")
		if authHeader == "" {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "missing authorization header",
			}).Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "invalid authorization scheme",
			}).Abort()
			return
		}

		if _, err := facades.Auth(ctx).Parse(token); err != nil {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "invalid or expired token",
			}).Abort()
			return
		}

		ctx.WithValue("token", token)
		ctx.Request().Next()
	}
}
