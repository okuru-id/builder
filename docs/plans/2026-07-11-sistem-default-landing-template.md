# Sistem Default Landing Template Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make supplied Sistem HTML the active default storefront landing and first admin gallery template, with iframe preview in admin gallery.

**Architecture:** Add nullable `html` column to landing_templates; branch Apply to write settings when template has HTML (no section replacement); seed Sistem template + activate it; render iframe preview in admin gallery when `html` present.

**Tech Stack:** Go 1.24, Goravel ORM, SQLite, Vue 3 + TypeScript, shadcn-vue Sheet, Tailwind v4.

Reference spec: `docs/superpowers/specs/2026-07-11-sistem-default-landing-template-design.md`

The exact HTML lives in `docs/plans/sistem-landing.html` (Task 0) and must be stored verbatim; do not change `#` links, fonts, or Tailwind CDN.

---

### Task 0: Save exact Sistem HTML reference file

**Files:**
- Create: `docs/plans/sistem-landing.html`

**Step 1: Create file with exact HTML supplied by user**

Write the full HTML document from the approved design into the file. Do not reformat, minify, or edit any links. This is the single source of truth for both seeder and verification.

**Step 2: Commit**

```bash
git add docs/plans/sistem-landing.html
git commit -m "docs: add Sistem landing HTML reference"
```

### Task 1: Add `html` column to landing_templates

**Files:**
- Create: `backend/database/migrations/20260712000001_add_html_to_landing_templates.go`
- Modify: `backend/bootstrap/migrations.go`

**Step 1: Create migration**

Create `backend/database/migrations/20260712000001_add_html_to_landing_templates.go`:

```go
package migrations

import (
	"okuru/app/facades"
)

type M20260712000001AddHtmlToLandingTemplates struct{}

func (r *M20260712000001AddHtmlToLandingTemplates) Signature() string {
	return "20260712000001_add_html_to_landing_templates"
}

func (r *M20260712000001AddHtmlToLandingTemplates) Up() error {
	if facades.Schema().HasTable("landing_templates") && !facades.Schema().HasColumns("landing_templates", "html") {
		return facades.Schema().Table("landing_templates", func(table schema.Blueprint) {
			table.Text("html").Nullable()
		})
	}
	return nil
}

func (r *M20260712000001AddHtmlToLandingTemplates) Down() error {
	if facades.Schema().HasTable("landing_templates") && facades.Schema().HasColumns("landing_templates", "html") {
		return facades.Schema().Table("landing_templates", func(table schema.Blueprint) {
			table.DropColumn("html")
		})
	}
	return nil
}
```

Import `schema` from `github.com/goravel/framework/contracts/database/schema`.

**Step 2: Register migration**

In `backend/bootstrap/migrations.go`, add `&migrations.M20260712000001AddHtmlToLandingTemplates{}` to the migration slice after the existing landing_templates migration registration. Follow exact registration pattern of existing migrations.

**Step 3: Compile**

Run: `cd backend && go build ./...`

Expected: PASS.

**Step 4: Commit**

```bash
git add backend/database/migrations/20260712000001_add_html_to_landing_templates.go backend/bootstrap/migrations.go
git commit -m "feat: add html column to landing_templates"
```

### Task 2: Expose `html` on model and controllers

**Files:**
- Modify: `backend/app/models/landing_template.go`
- Modify: `backend/app/http/controllers/admin/landing_template_controller.go`

**Step 1: Add model field**

In `backend/app/models/landing_template.go`, add to `LandingTemplate` struct after `Sections`:

```go
HTML string `gorm:"type:text" json:"html,omitempty"`
```

Keep existing field tags unchanged.

**Step 2: Accept `html` in Store and Update inputs**

In `landing_template_controller.go`, add `HTML string` (json `html`) to the local input struct used by both `Store` and `Update`. Copy the value to the model (`tmpl.HTML = input.HTML`) in both methods, matching the existing field copy pattern. Do not require it.

**Step 3: Compile**

Run: `cd backend && go build ./...`

Expected: PASS.

**Step 4: Commit**

```bash
git add backend/app/models/landing_template.go backend/app/http/controllers/admin/landing_template_controller.go
git commit -m "feat: expose landing template html field"
```

### Task 3: Branch Apply to support HTML templates

**Files:**
- Modify: `backend/app/http/controllers/admin/landing_template_controller.go`

**Step 1: Add HTML-apply branch**

At the top of `Apply`, after loading `tmpl`, insert:

```go
if tmpl.HTML != "" {
	if err := applyHTMLTemplate(tmpl.HTML); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"message": "Template applied"})
}
```

Keep the existing section-replacement logic untouched beneath this branch for section templates.

**Step 2: Add helper**

Add helper in same file:

```go
func applyHTMLTemplate(html string) error {
	settings := []models.Setting{
		{Key: "landing_mode", Value: "custom"},
		{Key: "landing_template_html", Value: html},
	}
	for _, s := range settings {
		var setting models.Setting
		if err := facades.Orm().Query().UpdateOrCreate(&setting, models.Setting{Key: s.Key}, s); err != nil {
			return err
		}
	}
	return nil
}
```

This writes both settings atomically per key and leaves existing landing sections unchanged.

**Step 3: Compile**

Run: `cd backend && go build ./...`

Expected: PASS.

**Step 4: Commit**

```bash
git add backend/app/http/controllers/admin/landing_template_controller.go
git commit -m "feat: apply html landing template writes settings"
```

### Task 4: Seed Sistem template and activate it

**Files:**
- Modify: `backend/database/seeders/landing_template_seeder.go`
- Modify: `backend/database/seeders/database_seeder.go`

**Step 1: Add Sistem as first template and make seeder idempotent**

In `landing_template_seeder.go`:

1. Prepend `models.LandingTemplate{Name: "Sistem", Description: "...", HTML: sistemHTML()}` as the first entry. Do not set Sections on it.
2. Add `func sistemHTML() string` returning the exact HTML from `docs/plans/sistem-landing.html` as a raw-string literal. Do not edit any byte.
3. Update the upsert loop to be idempotent on name: when `existing.ID != 0`, call `UpdateOrCreate` (or Save) on the struct using `Name` as the match attribute so the HTML stays current on re-seed. Keep creating when absent.

Description copy suggestion: `"Single-platform operations landing page. Tailwind + Inter/Space Grotesk."`

**Step 2: Activate Sistem on seed**

In `database_seeder.go` `seedSettings()` change the `landing_mode` and `landing_template_html` defaults to active Sistem:

- `landing_mode` value: `"custom"`
- `landing_template_html` value: `sistemHTML()` (call the new seeder-local helper). Keep `landing_template` as `default`.

Add a package-level `sistemHTML()` helper in `database_seeder.go` mirroring the seeder helper (or import from the seeder package) returning the same exact HTML; choose the option producing the smallest diff. Ensure no byte drift from Task 0 file.

**Step 3: Compile**

Run: `cd backend && go build ./...`

Expected: PASS.

**Step 4: Commit**

```bash
git add backend/database/seeders/landing_template_seeder.go backend/database/seeders/database_seeder.go
git commit -m "feat: seed and activate Sistem landing template"
```

### Task 5: Render iframe preview for HTML templates in admin gallery

**Files:**
- Modify: `frontend/src/views/LandingPage.vue`

**Step 1: Add `html` to template interface**

Extend local `LandingTemplate` interface with `html?: string`. Do not change other fields.

**Step 2: Add iframe preview helper**

Add:

```ts
function templatePreviewUrl(tmpl: LandingTemplate) {
  if (!tmpl.html) return ''
  const blob = new Blob([tmpl.html], { type: 'text/html' })
  return URL.createObjectURL(blob)
}
```

Track created blob URLs in `ref(new Map<number, string>())` and revoke them on unmount to avoid leaks.

**Step 3: Use iframe in card and Sheet when `html` present**

For both card media area and Sheet preview, when `tmpl.html` is truthy render an `<iframe>` sandboxed with `sandbox="allow-scripts"` (no `allow-same-origin`) using the blob URL. Apply classes so it scales inside the card media area (`pointer-events-none` and existing aspect classes) and fills the Sheet preview region.

When `html` is absent, keep existing image/fallback behavior unchanged.

**Step 4: Build**

Run: `cd frontend && bun run build`

Expected: PASS with only existing `@vueuse/core` warnings.

**Step 5: Commit**

```bash
git add frontend/src/views/LandingPage.vue
git commit -m "feat: preview html templates via iframe"
```

### Task 6: Verify end-to-end

**Files:** Modify only if defect found.

**Step 1: Backend build**

Run: `cd backend && go build ./...`

Expected: PASS.

**Step 2: Apply migration and seed fresh**

Run:

```bash
cd backend
./artisan migrate
./artisan db:seed
```

Expected: completes without error.

**Step 3: Manual assertions**

1. Visit storefront `/`; it renders Sistem HTML with Tailwind, fonts, and footer copy exactly.
2. Visit `/admin/landing-page`; Sistem is the first card and shows an iframe preview of the HTML.
3. Open Sistem Sheet; iframe preview fills the preview area; Cancel leaves sections and settings unchanged.
4. Apply Sistem via Sheet confirmation; verify `landing_mode=custom` and `landing_template_html` equal Sistem HTML, and existing landing sections are unchanged in `/admin/landing-page/content`.
5. Apply an existing section template (Professional); verify sections are replaced and `landing_mode` behavior reverts per existing section flow.

**Step 4: Repair commit if needed**

```bash
git add backend frontend
git commit -m "fix: complete Sistem landing flow"
```

Skip when no repair needed.
