# Phase 0: Cleanup Legacy Landing System

**Status:** ✅ done
**Commit:** `998200c`
**Goal:** Hapus seluruh legacy landing infra, verifikasi Postgres siap untuk builder baru.

## Checklist

- [x] Delete `backend/app/models/landing_section.go`
- [x] Delete `backend/app/models/landing_template.go`
- [x] Delete `backend/app/http/controllers/landing_controller.go`
- [x] Delete `backend/app/http/controllers/admin/landing_section_controller.go`
- [x] Delete `backend/app/http/controllers/admin/landing_template_controller.go`
- [x] Delete 4 legacy migrations (landing_sections + landing_templates)
- [x] Delete `backend/database/seeders/landing_section_seeder.go`
- [x] Delete `backend/database/seeders/landing_template_seeder.go`
- [x] Delete `frontend/src/views/LandingPage.vue`, `LandingPageContent.vue`, `LandingPagePreview.vue`
- [x] Delete `frontend/src/landing/` (landing SPA)
- [x] Hapus landing routes dari `api.go` + `web.go`
- [x] Hapus registrasi migration + seeder dari bootstrap
- [x] Hapus sidebar nav + router entries
- [x] Hapus settings `landing_mode`, `landing_template_html`
- [x] Ganti `frontend/index.html` dengan placeholder
- [x] Drop tables `landing_templates`, `landing_sections` (manual SQL)
- [x] Verifikasi `go build ./...` OK
- [x] Verifikasi Postgres live: DB `okuruid`, `DB_CONNECTION=postgres`

## Hasil

Server start clean, 0 reference ke landing lama. DB hanya berisi tabel non-landing. Storefront `/` = placeholder HTML sederhana.

## Catatan

- Konfigurasi sqlite di `config/database.go` dibiarkan (opsional, aman).
- Tabel legacy di-drop manual via `DROP TABLE` karena migration yang dihapus tidak otomatis drop.
