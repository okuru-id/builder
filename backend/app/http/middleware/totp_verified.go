package middleware

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

func TotpVerified() http.Middleware {
	return func(ctx http.Context) {
		userIDStr, err := facades.Auth(ctx).ID()
		if err != nil {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "unauthorized",
			}).Abort()
			return
		}

		userID, _ := strconv.ParseUint(userIDStr, 10, 64)

		var user models.User
		if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "user not found",
			}).Abort()
			return
		}

		if user.TotpSecret != nil && !user.TotpVerified {
			ctx.Response().Json(http.StatusForbidden, http.Json{
				"error":         "TOTP verification required",
				"totp_required": true,
			}).Abort()
			return
		}

		ctx.Request().Next()
	}
}
