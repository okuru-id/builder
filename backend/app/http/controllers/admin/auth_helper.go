package admin

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
)

// currentUserID returns the authenticated admin's ID, or 0 when unauthenticated.
// Every admin controller scopes its reads and writes through this so a logged-in
// user only ever sees rows they own.
func currentUserID(ctx http.Context) uint {
	idStr, err := facades.Auth(ctx).ID()
	if err != nil || idStr == "" {
		return 0
	}
	id, _ := strconv.ParseUint(idStr, 10, 64)
	return uint(id)
}
