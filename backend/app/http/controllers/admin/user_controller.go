package admin

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type userRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	IsActive *bool  `json:"is_active"`
	IsAdmin  *bool  `json:"is_admin"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	TotpOn   bool   `json:"totp_enabled"`
}

func toResponse(u models.User) userResponse {
	return userResponse{
		ID:       u.ID,
		Email:    u.Email,
		Name:     u.Name,
		IsActive: u.IsActive,
		IsAdmin:  u.IsAdmin,
		TotpOn:   u.TotpSecret != nil && u.TotpVerified,
	}
}

// isCurrentUserAdmin reports whether the caller has the is_admin flag.
// Admins manage all users; everyone else self-manages only.
func isCurrentUserAdmin(ctx http.Context) bool {
	uid := currentUserID(ctx)
	if uid == 0 {
		return false
	}
	var u models.User
	if err := facades.Orm().Query().Where("id = ?", uid).First(&u); err != nil || u.ID == 0 {
		return false
	}
	return u.IsAdmin
}

// keep old name as alias so other call sites compile
func isFirstAdmin(ctx http.Context) bool { return isCurrentUserAdmin(ctx) }

// Index: first admin sees all users; others see only themselves.
func (c *UserController) Index(ctx http.Context) http.Response {
	uid := currentUserID(ctx)
	if isFirstAdmin(ctx) {
		var users []models.User
		facades.Orm().Query().OrderBy("created_at", "desc").Get(&users)
		out := make([]userResponse, 0, len(users))
		for _, u := range users {
			out = append(out, toResponse(u))
		}
		return ctx.Response().Success().Json(http.Json{"data": out, "is_super": true})
	}
	var self models.User
	if err := facades.Orm().Query().Where("id = ?", uid).First(&self); err != nil || self.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": []userResponse{toResponse(self)}, "is_super": false})
}

// Show: first admin any user; others only self.
func (c *UserController) Show(ctx http.Context) http.Response {
	uid := currentUserID(ctx)
	targetID, _ := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if !isFirstAdmin(ctx) && uint(targetID) != uid {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "forbidden"})
	}
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", targetID).First(&user); err != nil || user.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"data": toResponse(user)})
}

// Store: first admin only.
func (c *UserController) Store(ctx http.Context) http.Response {
	if !isFirstAdmin(ctx) {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "forbidden"})
	}
	var input userRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "email and password are required"})
	}

	var existing models.User
	_ = facades.Orm().Query().Where("email = ?", input.Email).First(&existing)
	if existing.ID != 0 {
		return ctx.Response().Json(http.StatusConflict, http.Json{"error": "email already in use"})
	}

	hashed, err := facades.Hash().Make(input.Password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	user := models.User{
		Email:    input.Email,
		Name:     strings.TrimSpace(input.Name),
		Password: hashed,
		IsActive: input.IsActive == nil || *input.IsActive,
		IsAdmin:  input.IsAdmin != nil && *input.IsAdmin,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{"data": toResponse(user)})
}

// Update: first admin any user; others only self and only profile fields.
func (c *UserController) Update(ctx http.Context) http.Response {
	uid := currentUserID(ctx)
	targetID, _ := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	superAdmin := isFirstAdmin(ctx)
	if !superAdmin && uint(targetID) != uid {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "forbidden"})
	}

	var user models.User
	if err := facades.Orm().Query().Where("id = ?", targetID).First(&user); err != nil || user.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}

	var input userRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
	}

	updates := map[string]any{}
	if name := strings.TrimSpace(input.Name); name != "" {
		updates["name"] = name
	}
	if email := strings.TrimSpace(strings.ToLower(input.Email)); email != "" && email != user.Email {
		var clash models.User
		_ = facades.Orm().Query().Where("email = ?", email).First(&clash)
		if clash.ID != 0 && clash.ID != user.ID {
			return ctx.Response().Json(http.StatusConflict, http.Json{"error": "email already in use"})
		}
		updates["email"] = email
	}
	if input.Password != "" {
		hp, err := facades.Hash().Make(input.Password)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
		}
		updates["password"] = hp
	}
	// is_active only writable by admin
	if superAdmin && input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	// is_admin writable by admin, but never allow self-demotion (lockout guard)
	if superAdmin && input.IsAdmin != nil && uint(targetID) != uid {
		updates["is_admin"] = *input.IsAdmin
	}

	if len(updates) > 0 {
		if _, err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", user.ID).Update(updates); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
		}
	}
	// TOTP reset: super-admin only
	if superAdmin && ctx.Request().Query("reset_totp", "") == "1" {
		_, _ = facades.Orm().Query().Model(&models.User{}).Where("id = ?", user.ID).
			Update(map[string]any{"totp_secret": nil, "totp_verified": false})
	}

	_ = facades.Orm().Query().Where("id = ?", user.ID).First(&user)
	return ctx.Response().Success().Json(http.Json{"data": toResponse(user)})
}

// Destroy: super-admin only, never self.
func (c *UserController) Destroy(ctx http.Context) http.Response {
	if !isFirstAdmin(ctx) {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"error": "forbidden"})
	}
	id, _ := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if uint(id) == currentUserID(ctx) {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"error": "cannot delete your own account"})
	}
	res, err := facades.Orm().Query().Where("id = ?", id).Delete(&models.User{})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	if res.RowsAffected == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
	}
	return ctx.Response().Success().Json(http.Json{"deleted": true})
}
