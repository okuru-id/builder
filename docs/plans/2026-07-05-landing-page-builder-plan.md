# Landing Page Builder Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a dynamic section-based landing page builder with admin CRUD and public API.

**Architecture:** Single `landing_sections` table with typed JSON content. Admin Vue SPA consumes admin API for CRUD. Public API serves active sections to the landing page SPA.

**Tech Stack:** Goravel (Go), PostgreSQL JSONB, Vue 3 + shadcn-vue, TanStack Table

---

### Task 1: Backend — Migration + Model

**Files:**
- Create: `backend/database/migrations/20260705000001_create_landing_sections_table.go`
- Create: `backend/app/models/landing_section.go`
- Modify: `backend/bootstrap/migrations.go` — register migration

**Step 1: Create migration**

```go
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateLandingSectionsTable struct{}

func (m *CreateLandingSectionsTable) Signature() string {
	return "20260705000001_create_landing_sections_table"
}

func (m *CreateLandingSectionsTable) Up() error {
	return facades.Schema().Create("landing_sections", func(table schema.Blueprint) {
		table.ID()
		table.String("type").Unique()
		table.JSON("content")
		table.Integer("sort_order").Default(0)
		table.Bool("is_active").Default(true)
		table.Timestamps()
	})
}

func (m *CreateLandingSectionsTable) Down() error {
	return facades.Schema().DropIfExists("landing_sections")
}
```

**Step 2: Create model**

```go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/goravel/framework/database/orm"
)

type LandingSection struct {
	orm.Model
	Type      string          `gorm:"uniqueIndex" json:"type"`
	Content   LandingContent  `gorm:"type:jsonb" json:"content"`
	SortOrder int             `gorm:"default:0" json:"sort_order"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
}

type LandingContent map[string]any

func (c LandingContent) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *LandingContent) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}
```

**Step 3: Register migration**

In `backend/bootstrap/migrations.go`, add to the list pointer.

**Step 4: Run migration**

```bash
cd backend && ./artisan migrate
```

Expected: table `landing_sections` created.

**Step 5: Commit**

```bash
git add backend/database/migrations/20260705000001_create_landing_sections_table.go backend/app/models/landing_section.go backend/bootstrap/migrations.go
git commit -m "feat: add landing_sections migration and model"
```

---

### Task 2: Backend — Public Landing Controller

**Files:**
- Create: `backend/app/http/controllers/landing_controller.go`

**Step 1: Create controller**

```go
package controllers

import (
	"net/http"
	"okuru/app/facades"
	"okuru/app/models"
	"github.com/goravel/framework/contracts/http"
)

type LandingController struct{}

func NewLandingController() *LandingController {
	return &LandingController{}
}

func (c *LandingController) Index(ctx http.Context) http.Response {
	var sections []models.LandingSection
	err := facades.Orm().Query().
		Where("is_active", true).
		OrderBy("sort_order asc").
		Find(&sections)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": err.Error()})
	}

	result := make(map[string]any)
	for _, s := range sections {
		result[s.Type] = s.Content
	}

	return ctx.Response().Success().Json(http.Json{"data": result})
}
```

**Step 2: Commit**

```bash
git add backend/app/http/controllers/landing_controller.go
git commit -m "feat: add public landing page API controller"
```

---

### Task 3: Backend — Admin Landing Sections Controller

**Files:**
- Create: `backend/app/http/controllers/admin/landing_section_controller.go`

**Step 1: Create controller with CRUD methods**

Methods:
- `Index` — list all sections (admin needs inactive too)
- `Store` — create new section
- `Update` — update section content
- `Destroy` — soft/hard delete section
- `Toggle` — toggle is_active
- `Sort` — update sort_order

**Step 2: Commit**

```bash
git add backend/app/http/controllers/admin/landing_section_controller.go
git commit -m "feat: add admin landing sections CRUD controller"
```

---

### Task 4: Backend — Routes Registration

**Files:**
- Modify: `backend/routes/api.go`

**Step 1: Add public route**

```go
landingController := controllers.NewLandingController()
apiV1.Get("/landing", landingController.Index)
```

**Step 2: Add admin routes (inside JWT+TOTP group)**

```go
adminLandingCtrl := admin.NewLandingSectionController()
adminApi.Get("/landing-sections", adminLandingCtrl.Index)
adminApi.Post("/landing-sections", adminLandingCtrl.Store)
adminApi.Put("/landing-sections/{id}", adminLandingCtrl.Update)
adminApi.Delete("/landing-sections/{id}", adminLandingCtrl.Destroy)
adminApi.Patch("/landing-sections/{id}/toggle", adminLandingCtrl.Toggle)
adminApi.Patch("/landing-sections/{id}/sort", adminLandingCtrl.Sort)
```

**Step 3: Commit**

```bash
git add backend/routes/api.go
git commit -m "feat: register landing page routes"
```

---

### Task 5: Backend — Seeder

**Files:**
- Create: `backend/database/seeders/landing_section_seeder.go`
- Modify: `backend/bootstrap/seeders.go`

**Step 1: Create seeder with current hardcoded data**

Populate 5 sections: hero, clients, services, projects, cta — using data from `frontend/src/landing/App.vue`.

**Step 2: Commit**

```bash
git add backend/database/seeders/landing_section_seeder.go backend/bootstrap/seeders.go
git commit -m "feat: add landing section seeder"
```

---

### Task 6: Frontend — Admin Landing Page View

**Files:**
- Modify: `frontend/src/views/LandingPage.vue` (full rewrite)

**Step 1: Build admin UI**

- Fetch all sections from `GET /admin/api/landing-sections`
- Render each section as a card:
  - Section type as header
  - Summary preview (first item text, item count for lists)
  - Action buttons: Edit, Toggle active/inactive, Delete
- Edit opens a dialog with dynamic form based on section type:
  - Hero: text inputs for greeting, description, image paths, CTA
  - Clients: table with add/edit/delete rows (name, logo)
  - Services: table with add/edit/delete rows (title, description, icon)
  - Projects: table with add/edit/delete rows (title, description, github, technologies)
  - CTA: text inputs for heading, email, whatsapp
- Save per section via `PUT`
- Toggle active via `PATCH /{id}/toggle`

**Step 2: Commit**

```bash
git add frontend/src/views/LandingPage.vue
git commit -m "feat: rewrite landing page admin view with section builder"
```

---

### Task 7: Frontend — Landing Page SPA API Integration

**Files:**
- Modify: `frontend/src/landing/App.vue`

**Step 1: Fetch landing data on mount**

Replace hardcoded data with API call to `GET /api/v1/landing`. Fall back to hardcoded defaults if fetch fails.

**Step 2: Commit**

```bash
git add frontend/src/landing/App.vue
git commit -m "feat: integrate landing page with dynamic API data"
```

---

### Task 8: Verify

**Step 1: Run backend build**

```bash
cd backend && go build ./...
```

**Step 2: Run frontend build**

```bash
cd frontend && bun run build
```

**Step 3: If both pass, commit any fixes**

```bash
git commit -m "chore: fix build issues"
```
