package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
	"okuru/app/models"
)

type ShopController struct{}

func NewShopController() *ShopController {
	return &ShopController{}
}

// GET /shop — list active products.
func (c *ShopController) Index(ctx http.Context) http.Response {
	var products []models.Product
	facades.Orm().Query().Where("status = ?", "active").OrderBy("id", "desc").Get(&products)

	return ctx.Response().View().Make("shop/index.tmpl", map[string]any{
		"Title":    "Products",
		"Products": products,
	})
}

// GET /product/{slug} — single active product detail.
func (c *ShopController) Product(ctx http.Context) http.Response {
	slug := ctx.Request().Route("slug")

	var product models.Product
	if err := facades.Orm().Query().
		Where("slug = ?", slug).
		Where("status = ?", "active").
		First(&product); err != nil || product.ID == 0 {
		return ctx.Response().String(http.StatusNotFound, "Product not found")
	}

	return ctx.Response().View().Make("shop/product.tmpl", map[string]any{
		"Title":   product.Title,
		"Product": product,
	})
}

// POST /checkout/{slug} — payment integration placeholder (Phase 2).
func (c *ShopController) Checkout(ctx http.Context) http.Response {
	return ctx.Response().String(http.StatusNotImplemented, "Payment integration coming soon")
}
