# Landing Page Builder (Figma-like Canvas) Implementation Plan

> **Untuk Claude:** Gunakan superpowers:executing-plans untuk menjalankan plan ini task-per-task.

**Goal:** Replace section-picker form builder dengan canvas editor Figma-like (direct manipulation, tree JSON source of truth, codegen to HTML+Tailwind).

**Arsitektur:**
- DB Postgres JSONB (`pages.tree`) sebagai source of truth — bukan HTML string.
- Canvas Vue + `<NodeRenderer>` recursive render tree.
- Style = Tailwind utility classes (string array) per node → portable, codegen straightforward.
- Publish pipeline: tree → codegen → static HTML (cache, bukan source of truth).
- Replace total: landing_sections, landing_templates, landing_template_html dihapus.

**Tech Stack:** Go 1.24 + Goravel + GORM, PostgreSQL 17, Vue 3 + shadcn-vue + reka-ui + dnd-kit-vue + Pinia, Tailwind v4.

---

## Data Model (Phase 1 Foundation)

### Node Tree (JSONB structure)
```typescript
// Setiap node = satu entitas canvas (teks, frame, gambar, tombol, section, komponen)
interface Node {
  id: string                    // uuid
  type: 'frame' | 'text' | 'image' | 'button' | 'link' | 'section' | 'component'
  name: string                  // display name for tree panel
  props: Record<string, any>    // content: text, src, href, alt, placeholder, dll
  classes: string[]             // Tailwind utility classes — source of truth untuk style
  children: Node[]
  // Component instance (Phase 4)
  componentId?: number
  instanceOverrides?: Partial<Node>
}
```

### DB Tables
```
landing_pages
  id, slug UNIQUE, name, status (draft/published),
  tree JSONB (working draft), published_tree JSONB NULL (publish snapshot),
  version int DEFAULT 0, timestamps

landing_page_revisions     — append-only log
  id, landing_page_id FK, tree JSONB, message, created_at

landing_components         — reusable component master
  id, name, tree JSONB, timestamps
```

### Style = Tailwind classes
Tidak ada inline `style:` di tree. Semua style = Tailwind utility `classes[]`.
Contoh frame dengan auto-layout:
```json
"classes": ["flex", "flex-col", "gap-4", "p-6", "items-center", "max-w-6xl", "mx-auto"]
```
Codegen = `join(' ', node.classes)` + recursive render children.
Responsive override (Phase 3+) = array ekstra per breakpoint:
```json
"classes": ["flex", "gap-4", "p-4"],
"classesSm": ["flex-col"],
"classesLg": ["gap-6", "p-6"]
```

### API Endpoints (Backend)
```
GET    /admin/api/landing-pages               — list
POST   /admin/api/landing-pages               — create
GET    /admin/api/landing-pages/:id            — show (with tree)
PUT    /admin/api/landing-pages/:id            — autosave (update tree, create revision)
POST   /admin/api/landing-pages/:id/publish    — copy tree→published_tree, status=published
GET    /admin/api/landing-pages/:id/revisions  — list revisions
POST   /admin/api/landing-pages/:id/revisions/:rid/restore — rollback

GET    /admin/api/landing-pages/:id/components — list components on page
CRUD   /admin/api/landing-components           — component library
```

### Frontend Routes
```
/builder/:id    → CanvasEditor (full-screen, tanpa sidebar admin)
```

---

## Phase 0: Cleanup + PostgreSQL Consolidation

**Goal:** Hapus legacy landing system, verifikasi Postgres siap.

**Files:**
- Delete: `backend/app/models/landing_section.go`
- Delete: `backend/app/models/landing_template.go`
- Delete: `backend/app/http/controllers/landing_controller.go`
- Delete: `backend/app/http/controllers/admin/landing_section_controller.go`
- Delete: `backend/app/http/controllers/admin/landing_template_controller.go`
- Delete: `backend/database/migrations/20260705000001_create_landing_sections_table.go`
- Delete: `backend/database/migrations/20260706000001_drop_unique_index_on_landing_sections_type.go`
- Delete: `backend/database/migrations/20260711000001_create_landing_templates_table.go`
- Delete: `backend/database/migrations/20260712000001_add_html_to_landing_templates.go`
- Delete: `backend/database/seeders/landing_section_seeder.go`
- Delete: `backend/database/seeders/landing_template_seeder.go`
- Delete: `frontend/src/views/LandingPage.vue`
- Delete: `frontend/src/views/LandingPageContent.vue`
- Delete: `frontend/src/views/LandingPagePreview.vue`
- Delete: `frontend/src/views/landing/`` (directory, landing SPA)
- Modify: `backend/app/http/controllers/admin/` — hapus registrasi controller dari provider/dependencies
- Modify: `backend/routes/api.go` — hapus landing routes
- Modify: `backend/routes/web.go` — hapus custom HTML storefront render (`if landing_mode=custom` block)
- Modify: `backend/bootstrap/migrations.go` — hapus registrasi landing migrations
- Modify: `backend/bootstrap/seeders.go` — hapus ref landing seeder
- Modify (later): tambah registrasi landing_pages migration baru
- Modify: `backend/config/database.go` — hapus sqlite connection? (opsional, aman biarkan)

**Step 0.1: Hapus legacy landing files**

Run: `rm backend/app/models/landing_section.go` (dst untuk semua file di daftar)
Run: `rm -rf frontend/src/landing/`
Run: `rm -rf admin/` (jika belum cleanup)

**Step 0.2: Hapus landing routes dari api.go + web.go**

api.go: hapus route registrasi landingSectionController, landingTemplateController
web.go: hapus block `if landing_mode=custom` (ganti render default)

**Step 0.3: Hapus registrasi migration + seeder**

migrations.go: hapus import & baris landing migrations
seeders.go: hapus import & panggil landing seeder

**Step 0.4: Verifikasi build**

Run: `go build ./...` — harus OK
Run: `PGPASSWORD=password psql -h 127.0.0.1 -U postgres -d okuruid -c '\dt'` — landing_templates, landing_sections masih ada (data, tidak di-drop — migration yang dihapus tidak otomatis drop table. Tambah migration baru untuk DROP TABLE atau lakukan manual via SQL)

**Step 0.5: Drop legacy tables (via artisan atau migration baru)**

Opsi: migration `20260713000001_drop_landing_legacy.go` yang `DROP TABLE IF EXISTS landing_templates, landing_sections CASCADE`.

**Step 0.6: Commit**

```bash
git add -A
git commit -m "Phase 0: cleanup legacy landing system, ready for builder"
```

---

## Phase 1: Data Model — Tree JSON + API

**Goal:** landing_pages + landing_page_revisions + landing_components tables, models, CRUD API, autosave.

**Files:**
- Create: `backend/database/migrations/20260713000001_create_landing_pages_tables.go`
- Create: `backend/app/models/landing_page.go`
- Create: `backend/app/models/landing_page_revision.go`
- Create: `backend/app/models/landing_component.go`
- Create: `backend/app/http/controllers/admin/landing_page_controller.go`
- Modify: `backend/bootstrap/migrations.go` — register M20260713000001CreateLandingPagesTables
- Modify: `backend/routes/api.go` — register landing-pages routes
- Modify: `backend/routes/web.go` — simplify storefront `/`

### Models

`backend/app/models/landing_page.go`:
```go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/goravel/framework/database/orm"
)

type LandingPage struct {
	orm.Model
	Slug          string           `gorm:"uniqueIndex" json:"slug"`
	Name          string           `json:"name"`
	Status        string           `gorm:"default:draft" json:"status"`
	Tree          LandingPageTree  `gorm:"type:jsonb" json:"tree"`
	PublishedTree *LandingPageTree `gorm:"type:jsonb" json:"published_tree,omitempty"`
	Version       int              `gorm:"default:0" json:"version"`
}

// LandingPageTree = raw JSON node sebagai map
type LandingPageTree map[string]any

func (t LandingPageTree) Value() (driver.Value, error) {
	if t == nil { return nil, nil }
	return json.Marshal(t)
}
func (t *LandingPageTree) Scan(val any) error {
	if val == nil { *t = nil; return nil }
	var raw string
	switch v := val.(type) {
	case []byte: raw = string(v)
	case string: raw = v
	}
	return json.Unmarshal([]byte(raw), t)
}
```

`backend/app/models/landing_page_revision.go`:
```go
type LandingPageRevision struct {
	orm.Model
	LandingPageID uint            `gorm:"index" json:"landing_page_id"`
	Tree          LandingPageTree `gorm:"type:jsonb" json:"tree"`
	Message       string          `json:"message,omitempty"`
	CreatedAt     *carbon.DateTime `gorm:"autoCreateTime" json:"created_at"`
}
```

`backend/app/models/landing_component.go`:
```go
type LandingComponent struct {
	orm.Model
	Name string          `json:"name"`
	Tree LandingPageTree `gorm:"type:jsonb" json:"tree"`
}
```

### Controller: LandingPageController

Methods: Index, Show, Store, Update, Publish, Revisions, RestoreRevision

Flow:
- `Store`: generate uuid slug dari name, tree = default root frame
- `Update`: update tree, increment version, CREATE revision record
- `Publish`: copy tree → published_tree, set status=published
- `Revisions`: list revisions DESC
- `RestoreRevision`: replacement rev tree → current tree, tambah version

### Routes (register in api.go)

```go
r.Get("/landing-pages", landingPageController.Index)
r.Post("/landing-pages", landingPageController.Store)
r.Get("/landing-pages/{id}", landingPageController.Show)
r.Put("/landing-pages/{id}", landingPageController.Update)
r.Post("/landing-pages/{id}/publish", landingPageController.Publish)
r.Get("/landing-pages/{id}/revisions", landingPageController.Revisions)
r.Post("/landing-pages/{id}/revisions/{rid}/restore", landingPageController.RestoreRevision)
```

### Storefront / (web.go)

Ganti blok kompleks dengan render sederhana:
- Cari landing_page where status=published
- Jika ada → render tree ke HTML (placeholder, pakai template minimal atau codegen nanti)
- Jika tidak → `return ctx.Response().String(200, "<html><body><h1>Landing page not published</h1></body></html>")`

Atau simpan `published_html` di DB (goceng hasil codegen di publish) → serve langsung. Tambah field: `published_html text`.

**Fase awal**: storefront serve placeholder text. Lanjut Phase 5 untuk codegen.

### Commit Phase 1

```bash
git add backend/database/migrations/20260713000001_create_landing_pages_tables.go
git add backend/app/models/landing_page.go
git add backend/app/models/landing_page_revision.go
git add backend/app/models/landing_component.go
git add backend/app/http/controllers/admin/landing_page_controller.go
git add backend/bootstrap/migrations.go
git add backend/routes/api.go backend/routes/web.go
git commit -m "Phase 1: landing pages tree JSON data model + CRUD API"
```

---

## Phase 2: Backend Publish + Codegen Foundation

**Goal:** Publish endpoint runs tree → HTML codegen, stores `published_html` on page. Storefront serves it.

**Files:**
- Create: `backend/app/services/codegen.go` — walk tree JSON → HTML string
- Modify: `backend/app/http/controllers/admin/landing_page_controller.go` — Publish runs codegen
- Modify: `backend/app/models/landing_page.go` — tambah `PublishedHTML string`
- Modify: `backend/routes/web.go` — serve `published_html`

### Codegen Services (codegen.go)

```go
package services

func GenerateHTML(tree map[string]any, pageName string) string {
	// Walk tree recursively
	// Each node: type → HTML element, props → attributes, classes → class=""
	// text → <span class="classes">{props.text}</span>
	// frame → <div class="classes">{children}</div>
	// image → <img class="classes" src={props.src} alt={props.alt} />
	// button → <button class="classes">{props.text}</button>
	// section → <section class="classes">{children}</section>
	// Wrap dengan <!DOCTYPE html><html>... Tailwind CDN script
	return html
}
```

### Commit Phase 2

```bash
git add backend/app/services/codegen.go
git add -A
git commit -m "Phase 2: publish pipeline, tree→HTML codegen, storefront render"
```

---

## Phase 3: Frontend Builder — Canvas Editor

**Goal:** Builder route, recursive NodeRenderer, selection + inline edit, autosave.

**Files:**
- Create: `frontend/src/types/page-builder.ts` — TypeScript node types
- Modify: `frontend/src/router/index.ts` — tambah `/builder/:id` route (top-level, tanpa sidebar)
- Create: `frontend/src/views/Builder.vue` — builder layout (toolbar, canvas, panels)
- Create: `frontend/src/components/builder/NodeRenderer.vue` — recursive node renderer
- Create: `frontend/src/components/builder/Toolbar.vue` — publish, undo/redo, breakpoints
- Create: `frontend/src/components/builder/NodeTreePanel.vue` — tree outline sidebar (optional for MVP)
- Create: `frontend/src/components/builder/InspectorPanel.vue` — properties editor (right panel)
- Create: `frontend/src/components/builder/Canvas.vue` — canvas viewport, handles zoom/pan

### TypeScript Types

```typescript
// frontend/src/types/page-builder.ts
export interface Node {
  id: string
  type: 'frame' | 'text' | 'image' | 'button' | 'link' | 'section' | 'component'
  name: string
  props: Record<string, any>
  classes: string[]
  children: Node[]
  componentId?: number
  instanceOverrides?: Partial<Node>
}

export interface Page {
  id: number
  slug: string
  name: string
  status: 'draft' | 'published'
  tree: { root: Node }
  version: number
  created_at: string
  updated_at: string
}
```

### NodeRenderer.vue (core)
```vue
<script setup lang="ts">
import type { Node } from '@/types/page-builder'

const props = defineProps<{ node: Node; depth?: number }>()
const emit = defineEmits<{ select: [id: string] }>()
</script>

<template>
  <component :is="node.type === 'text' ? 'span' :
                   node.type === 'heading' ? ['h1','h2','h3','h4'][node.props.level || 0] || 'h2' :
                   node.type === 'image' ? 'img' :
                   node.type === 'button' ? 'button' :
                   node.type === 'link' ? 'a' :
                   node.type === 'section' ? 'section' :
                   'div'"
    :class="node.classes"
    :href="node.props.href"
    :src="node.props.src"
    :alt="node.props.alt"
    @click.stop="$emit('select', node.id)"
  >
    <template v-if="node.type === 'text' || node.type === 'heading'">
      {{ node.props.text }}
    </template>
    <template v-else-if="node.type === 'image'" />
    <template v-else>
      <NodeRenderer v-for="child in node.children" :key="child.id"
        :node="child" :depth="(depth ?? 0) + 1"
        @select="$emit('select', $event)" />
    </template>
  </component>
</template>
```

### Autosave Logic
```typescript
import { watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
const save = useDebounceFn(async () => {
  await api.put(`/landing-pages/${pageId}`, { tree })
}, 2000)

watch(() => tree.value, save, { deep: true })
```

### Inline Text Editing (contenteditable)
```vue
<span
  :contenteditable="isEditing"
  @blur="updateText($event.target.innerText)"
  @dblclick="startEditing"
  v-text="node.props.text"
/>
```

### Commit Phase 3

```bash
git add frontend/src/types/ frontend/src/views/Builder.vue frontend/src/components/builder/
git add frontend/src/router/index.ts
git add frontend/package.json  # (if pinia added)
git commit -m "Phase 3: canvas builder frontend — NodeRenderer, selection, inline edit, autosave"
```

---

## Phase 4: Advanced Canvas — Drag, Resize, Snap, Multi-select

**Goal:** Direct manipulation: drag move/reorder (dnd-kit-vue), multi-select, resize handles, smart snapping grid.

**Tasks:**
4.1: Integrate dnd-kit-vue untuk drag reorder di tree (child move antar sibling)
4.2: Selection overlay — highlight selected node border + handles
4.3: Multi-selection (Shift+click, lasso select)
4.4: Resize handles (frame corners, right/bottom edges)
4.5: Smart snapping — alignment guides saat drag (align top/left/center/right/bottom of nearest siblings)
4.6: Arrow key nudge selection (1px / 10px with Shift)
4.7: Cmd+D duplicate selected
4.8: Delete selected (Backspace/Delete)
4.9: Quick add node (click "+" on canvas edge / drag from palette)

**Commit:** `git commit -m "Phase 4: drag, resize, snap, multi-select, keyboard shortcuts"`

---

## Phase 5: Component/Instance System

**Goal:** Simpan reusable component, pasang instance, edit master → update semua instance.

**Tasks:**
5.1: CRUD `/landing-components` (backend)
5.2: "Create Component" — save selected node(s) as reusable
5.3: Component palette (left drawer, show available components)
5.4: Instance rendering — NodeRenderer resolve `componentId` → load master tree + merge `instanceOverrides`
5.5: Edit master component → flag instance dirty, show "Update from master" badge
5.6: Break link — convert instance to standalone copy

**Commit:** `git commit -m "Phase 5: component/instance system"`

---

## Phase 6: Style Panel & Token System

**Goal:** Right-side inspector dengan UI controls untuk Tailwind classes. Bukan raw textarea class.

**Tasks:**
6.1: Layout section — direction (flex-row/flex-col), align, justify, gap, padding, margin → dropdown + slider
6.2: Typography section — font family, size, weight, color → select + color picker
6.3: Spacing section — margin, padding per direction
6.4: Background — color, image
6.5: Border — width, color, radius
6.6: Width/Height — min, max, fixed
6.7: Token system — define brand tokens (primary, secondary, font-heading, spacing-unit) → apply as Tailwind classes
6.8: Each control mutates `node.classes[]` → rebuild

**Commit:** `git commit -m "Phase 6: style panel + token system"`

---

## Phase 7: One-Way Codegen (Tree → Vue/HTML)

**Goal:** Generate clean, downloadable Vue SFC or HTML+Tailwind dari tree.

**Tasks:**
7.1: Go codegen `services/codegen.go` — walk tree → HTML string (deterministic, idempotent)
7.2: TypeScript codegen mirror (untuk live preview export)
7.3: Download button — export HTML or Vue SFC
7.4: Codegen tests — input tree X → expected output hash
7.5: Safe list generator — scan tree classes → generate safelist.txt untuk Tailwind JIT build

**Tests (critical):**
```go
func TestGenerateHTML_EmptyTree(t *testing.T) { /* ... */ }
func TestGenerateHTML_Deterministic(t *testing.T) {
    html1 := GenerateHTML(tree)
    html2 := GenerateHTML(tree)
    assert.Equal(t, html1, html2)
}
func TestGenerateHTML_NestedFrame(t *testing.T) { /* ... */ }
```

**Commit:** `git commit -m "Phase 7: one-way codegen tree→HTML/Vue + tests"`

---

## Phase 8: Two-Way Sync (Code → Tree)

**Goal:** Parse existing `.vue`/HTML back to tree JSON — developer edit code, canvas reflect.

**Tasks:**
8.1: Parse Vue SFC via `@vue/compiler-sfc` → extract template AST
8.2: Convert AST nodes (v-for, v-if, v-bind) → builder node tree
8.3: Parse Tailwind classes back to inspector properties (reverse mapping)
8.4: Round-trip tests: generate → parse → generate = identical

**Commit:** `git commit -m "Phase 8: two-way sync code→tree"`

---

## Anti-Pattern Checklist (Brief §4)

| Check | Rule |
|-------|------|
| ✅ | Source of truth = tree JSON, bukan HTML string |
| ✅ | Canvas direct manipulation, bukan form-based-only |
| ✅ | Section library opsional (static starting point), bukan satu-satunya cara |
| ✅ | Token + Tailwind utilities, bukan inline style |
| ✅ | Publish pipeline terpisah dari canvas editing |
| ✅ | Output (codegen) pure Tailwind, 0 import `@/components/ui/*` |

---

## Command Cheat Sheet

```bash
# Backend
cd backend && go build ./...
CGO_ENABLED=0 go build -o okuru .
./artisan migrate
./artisan db:seed

# Frontend
cd frontend && bun install && bun run build
bun run dev

# DB
docker exec database-pgsql-1 psql -U postgres -d okuruid -c '\dt'
PGPASSWORD=password psql -h 127.0.0.1 -U postgres -d okuruid -c 'SELECT * FROM landing_pages;'

# Git
git log --oneline -10
git status --short
```
