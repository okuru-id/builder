package routes

import (
	"encoding/json"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"okuru/app/http/controllers"
	"okuru/app/facades"
	"okuru/app/models"
	"okuru/app/services"
)

func Web() {
	// Storefront root + preview. Path-based pages (/promo) are handled by Fallback.
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		html, status := resolveStorefront(ctx, "")
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(status, html)
	})

	// Admin SPA: serve the shell for /admin and any /admin/* client route.
	facades.Route().Get("/admin/", func(ctx http.Context) http.Response {
		return ctx.Response().File("./public/admin.html")
	})
	// SPA fallback + storefront path dispatch:
	//   - /admin/*           → admin shell (client router)
	//   - any other segment   → resolve as a published page by path
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		path := strings.TrimPrefix(ctx.Request().Path(), "/")
		if strings.HasPrefix(path, "admin") {
			return ctx.Response().File("./public/admin.html")
		}
		html, status := resolveStorefront(ctx, path)
		return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(status, html)
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

// resolveStorefront picks which published HTML to serve for a request.
//
// Resolution priority:
//  1. ?preview=<id>  → render the draft tree of that page (no publish required).
//  2. Host header     → page whose Domain matches (published only).
//  3. path == ""      → page flagged is_home (published), else the app's own
//     marketing landing (defaultStorefrontHTML).
//  4. path != ""      → page whose Path equals the segment (published).
//     Unknown path → 404.
//
// Returns (body, status). Callers wrap with Content-Type.
func resolveStorefront(ctx http.Context, path string) (string, int) {
	// 1. Preview mode (draft render). Allowed on any path/host.
	if previewID := ctx.Request().Query("preview", ""); previewID != "" {
		var page models.LandingPage
		if err := facades.Orm().Query().Where("id = ?", previewID).First(&page); err != nil || page.ID == 0 {
			return "Halaman tidak ditemukan", http.StatusNotFound
		}
		return renderPreviewHTML(page, ctx.Request().Host()), http.StatusOK
	}

	host := normalizeHost(ctx.Request().Host())

	// 2. Custom-domain match.
	if host != "" {
		var page models.LandingPage
		_ = facades.Orm().Query().
			Where("domain = ?", host).
			Where("status = ?", "published").
			First(&page)
		if page.ID != 0 && page.PublishedHTML != "" {
			return page.PublishedHTML, http.StatusOK
		}
	}

	// 3. Home.
	if path == "" {
		var page models.LandingPage
		_ = facades.Orm().Query().
			Where("is_home = ?", true).
			Where("status = ?", "published").
			First(&page)
		if page.ID != 0 && page.PublishedHTML != "" {
			return page.PublishedHTML, http.StatusOK
		}
		return defaultStorefrontHTML(ctx.Request().Host()), http.StatusOK
	}

	// 4. Path match.
	var page models.LandingPage
	_ = facades.Orm().Query().
		Where("path = ?", path).
		Where("status = ?", "published").
		First(&page)
	if page.ID != 0 && page.PublishedHTML != "" {
		return page.PublishedHTML, http.StatusOK
	}
	return "404 — halaman tidak ditemukan", http.StatusNotFound
}

// normalizeHost lowercases the Host header and strips the port.
func normalizeHost(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	if i := strings.LastIndex(h, ":"); i != -1 && !strings.Contains(h, "[") {
		h = h[:i]
	}
	return strings.Trim(h, ".")
}

// renderPreviewHTML converts a page's working tree into a full HTML document
// using the codegen pipeline. Used for ?preview= mode.
func renderPreviewHTML(page models.LandingPage, host string) string {
	var m map[string]any
	if err := json.Unmarshal(page.Tree, &m); err != nil {
		return defaultStorefrontHTML(host)
	}
	resolved := services.ResolveComponentInstances(m)
	return services.NewLandingCodegen().Generate(resolved, page.Name+" (preview)")
}

// defaultStorefrontHTML is the builder app's own marketing landing, shown at the
// apex domain when no page is flagged is_home + published. It advertises the
// builder itself so first-time visitors understand what the product is.
func defaultStorefrontHTML(host string) string {
	name := host
	if i := strings.LastIndex(host, ":"); i != -1 {
		name = host[:i]
	}
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>` + name + ` — Website & Landing Page Builder</title>
<meta name="description" content="A visual builder for landing pages. Drag, edit, and publish in minutes — no code required." />
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
<style>
  @keyframes floaty { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-8px)} }
  .floaty{animation:floaty 6s ease-in-out infinite}
  .grid-bg{background-image:linear-gradient(rgba(99,102,241,.07) 1px,transparent 1px),linear-gradient(90deg,rgba(99,102,241,.07) 1px,transparent 1px);background-size:32px 32px}
</style>
</head>
<body class="min-h-screen bg-neutral-950 text-neutral-100 antialiased">
  <div class="grid-bg min-h-screen">
    <!-- Nav -->
    <header class="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
      <div class="flex items-center gap-2">
        <span class="flex size-8 items-center justify-center rounded-lg bg-indigo-500 font-bold text-white">o</span>
        <span class="text-lg font-semibold tracking-tight">` + name + `</span>
      </div>
      <nav class="flex items-center gap-3 text-sm">
        <a href="/admin/" class="rounded-md bg-white/10 px-3 py-1.5 font-medium text-white transition hover:bg-white/20">Open Builder</a>
      </nav>
    </header>

    <!-- Hero -->
    <main class="mx-auto max-w-4xl px-6 pb-24 pt-16 text-center sm:pt-24">
      <span class="inline-flex items-center gap-1.5 rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-indigo-300">
        Visual · Drag-and-drop · AI-assisted
      </span>
      <h1 class="mx-auto mt-6 max-w-3xl text-4xl font-bold leading-tight tracking-tight sm:text-6xl">
        Build &amp; publish landing pages
        <span class="bg-gradient-to-r from-indigo-400 to-fuchsia-400 bg-clip-text text-transparent">in minutes</span>
      </h1>
      <p class="mx-auto mt-6 max-w-xl text-base text-neutral-400 sm:text-lg">
        A no-code website builder with a live canvas, reusable components, and one-click publishing to your own domain.
      </p>
      <div class="mt-9 flex items-center justify-center gap-3">
        <a href="/admin/" class="rounded-lg bg-indigo-500 px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-indigo-500/30 transition hover:bg-indigo-400">
          Start building
        </a>
        <a href="/admin/#/pages" class="rounded-lg border border-white/15 px-5 py-2.5 text-sm font-medium text-neutral-200 transition hover:bg-white/5">
          View dashboard
        </a>
      </div>

      <!-- Mock canvas preview -->
      <div class="floaty mx-auto mt-16 max-w-3xl overflow-hidden rounded-xl border border-white/10 bg-neutral-900 shadow-2xl">
        <div class="flex items-center gap-1.5 border-b border-white/10 px-4 py-2.5">
          <span class="size-2.5 rounded-full bg-red-400/80"></span>
          <span class="size-2.5 rounded-full bg-amber-400/80"></span>
          <span class="size-2.5 rounded-full bg-emerald-400/80"></span>
          <span class="ml-3 text-xs text-neutral-500">builder · canvas</span>
        </div>
        <div class="grid grid-cols-12 gap-2 p-4 text-left">
          <div class="col-span-3 space-y-2">
            <div class="h-4 rounded bg-white/5"></div>
            <div class="h-4 w-3/4 rounded bg-white/5"></div>
            <div class="h-4 w-1/2 rounded bg-indigo-500/40"></div>
          </div>
          <div class="col-span-9 space-y-2">
            <div class="h-24 rounded-lg bg-gradient-to-br from-indigo-500/20 to-fuchsia-500/20"></div>
            <div class="grid grid-cols-3 gap-2">
              <div class="h-16 rounded-lg bg-white/5"></div>
              <div class="h-16 rounded-lg bg-white/5"></div>
              <div class="h-16 rounded-lg bg-white/5"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Features -->
      <div class="mx-auto mt-20 grid max-w-3xl gap-8 sm:grid-cols-3">
        <div class="text-center">
          <div class="mx-auto flex size-10 items-center justify-center rounded-lg bg-indigo-500/15 text-indigo-300">🎨</div>
          <h3 class="mt-3 text-sm font-semibold">Live canvas</h3>
          <p class="mt-1 text-xs text-neutral-400">Edit visually. What you see is what you ship.</p>
        </div>
        <div class="text-center">
          <div class="mx-auto flex size-10 items-center justify-center rounded-lg bg-fuchsia-500/15 text-fuchsia-300">⚡</div>
          <h3 class="mt-3 text-sm font-semibold">Instant publish</h3>
          <p class="mt-1 text-xs text-neutral-400">One click goes live on your path or domain.</p>
        </div>
        <div class="text-center">
          <div class="mx-auto flex size-10 items-center justify-center rounded-lg bg-emerald-500/15 text-emerald-300">🧩</div>
          <h3 class="mt-3 text-sm font-semibold">Components</h3>
          <p class="mt-1 text-xs text-neutral-400">Reuse blocks across every page.</p>
        </div>
      </div>
    </main>

    <footer class="border-t border-white/10 py-6 text-center text-xs text-neutral-500">
      Powered by <a href="/admin/" class="text-neutral-300 hover:text-white">` + name + ` Builder</a>
    </footer>
  </div>
</body>
</html>`
}
