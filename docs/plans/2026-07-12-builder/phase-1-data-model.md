# Phase 1: Data Model — Tree JSON + CRUD API

**Status:** ✅ done
**Commit:** `d8bc2f2`
**Goal:** Tabel `landing_pages`, `landing_page_revisions`, `landing_components` (JSONB). Models + controller + routes. Autosave + revision log.

## Checklist

- [x] Migration `20260713000001_create_landing_pages_tables.go`
- [x] Model `landing_page.go` (Tree, PublishedTree JSONB, Version, Status)
- [x] Model `landing_page_revision.go` (append-only snapshot)
- [x] Model `landing_component.go` (reusable master)
- [x] Controller `landing_page_controller.go`: Index, Show, Store, Update, Publish, Revisions, RestoreRevision
- [x] Register migration di `bootstrap/migrations.go`
- [x] Register routes di `routes/api.go` (admin/api group, authenticated)
- [x] JSONB type = `gorm.io/datatypes.JSON` (native, handle marshal otomatis)
- [x] API verified end-to-end: create/update/revisions/publish/restore semua 200

## API Endpoints

```
GET    /admin/api/landing-pages                     → list (tanpa tree payload)
POST   /admin/api/landing-pages                     → create draft
GET    /admin/api/landing-pages/{id}                → show full tree
PUT    /admin/api/landing-pages/{id}                → autosave tree + revision
POST   /admin/api/landing-pages/{id}/publish        → snapshot tree → published_tree
GET    /admin/api/landing-pages/{id}/revisions      → list revisions DESC
POST   /admin/api/landing-pages/{id}/revisions/{rid}/restore  → rollback
```

## Default Tree (saat create)

```json
{
  "root": {
    "id": "root",
    "type": "frame",
    "name": "Page",
    "props": {},
    "classes": ["min-h-screen", "bg-white", "text-neutral-900"],
    "children": []
  }
}
```

## Hasil Verifikasi

```
login 200, create 200, update 200, revisions 200, publish 200, restore 200
```

## Catatan

- Awalnya pakai custom `LandingPageTree map[string]any` dengan Valuer/Scanner → GORM pg driver tidak handle, error `invalid input syntax for type json`. Fix: ganti ke `datatypes.JSON` (sudah ada sebagai indirect dep).
- Revision table awal tidak punya `updated_at` (orm.Model embed butuh keduanya). Fix: migration pakai `table.Timestamps()`.
