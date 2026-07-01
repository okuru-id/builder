# Design Doc: Okuru.id Admin Panel & Backend API

**Date:** 2026-07-02
**Status:** Approved
**Author:** Kurob

## Ringkasan

Backend API + admin panel untuk personal website okuru.id. Frontend React statis tetap di GitHub Pages. Backend dibangun dengan Goravel (Go) + SQLite, di-deploy via Docker ke VPS 512MB. Admin panel menggunakan Vue 3 + shadcn-vue. Storefront menggunakan Go template + Alpine.js.

## Konteks & Masalah

Saat ini okuru.id adalah static React site di GitHub Pages dengan konten hardcoded di 3 file JSON (`index-content.json`, `blog-content.json`, `content.json`). Tidak ada cara untuk update konten tanpa edit file JSON manual + redeploy. Contact form tidak berfungsi (placeholder). Blog posts belum punya halaman detail.

## Tujuan

1. CRUD konten via admin panel (blog, projects, open source, settings)
2. Contact form berfungsi (simpan ke database, inbox di admin)
3. Storefront `shop.okuru.id` untuk digital products + donasi
4. Auth aman (password + 2FA TOTP)
5. Resource efisien untuk VPS 512MB

## Tech Stack

| Komponen | Teknologi | Alasan |
|---|---|---|
| Backend | Goravel v1.17 (Go 1.23+) | Laravel-style API, single binary, efisien |
| Database | SQLite (WAL mode) | Zero overhead, file-based, cukup untuk personal site |
| Admin Panel | Vue 3 + Vite + shadcn-vue (Reka UI) | Request user, familiar pattern |
| Storefront | Go html/template + Alpine.js | No build step, SEO friendly, paling ringan |
| Reverse Proxy | Caddy 2 (Docker) | Auto HTTPS, config minimal |
| Deployment | Docker Compose + GitLab CI | Build di runner, VPS hanya pull |

## Arsitektur

```
okuru.id (GitHub Pages, React static)
  │ fetch JSON (GET only)
  ▼
api.okuru.id (VPS)
  ├── Caddy (auto-HTTPS, reverse proxy)
  ├── Goravel binary (Docker container)
  │   ├── /v1/*        Public API (CORS)
  │   ├── /admin/api/* Protected (JWT + TOTP)
  │   └── /admin/*     Vue SPA (embed.FS)
  └── SQLite (volume mount)

shop.okuru.id (VPS, same Goravel binary)
  └── Go template + Alpine.js storefront
```

Satu binary Goravel serve semua domain. Caddy route berdasarkan Host header.

## Data Model

### users
- id, email, password (bcrypt), totp_secret (encrypted), totp_verified, created_at

### posts
- id, slug (unique), title_en, title_id, excerpt_en, excerpt_id
- content_en, content_id (markdown), category (FK), tags (JSON)
- thumbnail, status (draft|published), published_at, read_time (auto-calc)
- created_at, updated_at

### categories
- id, slug (unique), name_en, name_id

### projects
- id, sort_order, title_en, title_id, description_en, description_id
- tech_stack (JSON), url, featured

### open_source_projects
- id, sort_order, title_en, title_id, description_en, description_id
- github_url, technologies (JSON), stars, license

### products (digital products)
- id, slug (unique), title, description, price (IDR integer)
- type (ebook|template|source_code|donation), file_path, thumbnail
- status (active|inactive), created_at, updated_at

### orders
- id, reference (unique), product_id (FK), buyer_email, buyer_name
- amount (IDR), status (pending|paid|failed|delivered)
- payment_ref, download_token (hashed), download_count, max_downloads
- expires_at, paid_at, created_at

### messages (contact inbox)
- id, name, email, message, status (unread|read|archived), created_at

### settings (key-value)
- id, key (unique), value (JSON)
- Keys: hero_title_en/id, hero_desc_en/id, start_year, social_links,
  skills, profile_image, cv_file

## API Routes

### Public API (api.okuru.id/v1) — CORS okuru.id + shop.okuru.id
- GET `/v1/content` — full site content (batch, untuk React load sekali)
- GET `/v1/posts` — list published (?category, ?page)
- GET `/v1/posts/:slug` — post detail
- GET `/v1/projects` — featured projects
- GET `/v1/open-source` — OS projects
- GET `/v1/skills` — tech stack
- GET `/v1/settings/public` — public settings
- POST `/v1/contact` — submit contact form (rate limited)

### Shop (shop.okuru.id) — server-rendered
- GET `/` — product listing
- GET `/product/:slug` — product detail
- POST `/checkout/:slug` — initiate Sumopod payment
- GET `/success?ref=xxx` — payment success
- GET `/download/:token` — download digital product
- POST `/webhook/sumopod` — payment webhook (verify signature)

### Admin API (api.okuru.id/admin/api) — JWT + TOTP
- POST `/admin/api/auth/login` — email+password → temp_token
- POST `/admin/api/auth/totp` — verify TOTP → JWT + refresh
- POST `/admin/api/auth/totp/setup` — generate secret + QR
- GET `/admin/api/auth/me` — current user
- CRUD `/admin/api/posts`
- CRUD `/admin/api/projects`
- CRUD `/admin/api/open-source`
- CRUD `/admin/api/categories`
- CRUD `/admin/api/products`
- GET/PUT `/admin/api/orders`
- GET/POST/DELETE `/admin/api/messages`
- GET/PUT `/admin/api/settings`
- POST `/admin/api/upload`

## Auth & Security

### Two-phase Login (Password + TOTP)
1. Phase 1: Verify email+password → return `temp_token` (2 min TTL, scope: totp_verify)
2. Phase 2: Verify TOTP code → return `access_token` (JWT 15 min) + `refresh_token` (7 day, HttpOnly cookie)

### Security Measures
- Password: bcrypt (Goravel `facades.Hash()`)
- TOTP secret: AES-256-GCM encrypted at rest (`APP_KEY`), library `pquerna/otp`
- JWT: Goravel built-in, short TTL, refresh via HttpOnly cookie
- HTTPS: Caddy auto Let's Encrypt + HSTS
- CORS: whitelist `okuru.id`, `shop.okuru.id`
- Rate limit: login 5/5min, contact 3/hour, webhook 30/min
- CSRF: same-origin admin SPA, token untuk mutations
- SQL injection: ORM prepared statements
- XSS: goldmark markdown dengan strict HTML policy
- File upload: MIME validation, max 5MB, random filename, outside webroot
- Download token: SHA-256 hashed, one-time, expire 72h, max 3 downloads
- Secrets: `.env` gitignored, `./artisan env:encrypt` untuk commit

### Admin Middleware
JWT claim `totp_v: true/false`. Middleware `RequireTotp` cek claim. Jika belum verified, hanya akses `/auth/totp/*`.

## Deployment (Docker + GitLab CI)

### Build Pipeline
Build di GitLab Runner (shared), push ke GitLab Container Registry. VPS hanya pull image + restart.

```
push → lint → test → build:image → publish → deploy
                                          ├─ staging (auto, develop branch)
                                          └─ production (manual, main branch)
```

### Docker Compose (Production VPS)
- app: image dari registry, memory limit 256M, healthcheck
- caddy: 2-alpine, memory limit 64M, auto-HTTPS
- volumes: db-data, storage-data, caddy-data, caddy-config
- network: bridge

### Dockerfile (Multi-stage)
1. Stage 1: node:20-alpine → build Vue admin (pnpm)
2. Stage 2: golang:1.23-alpine → build Go backend (CGO untuk SQLite)
3. Stage 3: alpine:3.20 → runtime (~25-30MB image)

### Resource Budget (512MB VPS)
```
OS              ~80 MB
Docker daemon   ~30 MB
Caddy           ~20 MB
Goravel (peak) ~120 MB
SQLite          ~5 MB
─────────────────────────
Total          ~290 MB
Headroom       ~222 MB
```

### GitLab CI Variables
- VPS_HOST, VPS_USER, VPS_SSH_KEY (File), VPS_SSH_FINGERPRINT, GITLAB_DEPLOY_TOKEN

### Git Workflow
```
feature/* → MR (lint+test)
develop   → staging (auto deploy)
main      → production (manual deploy)
v*.*.*    → release tag (manual deploy)
```

### SQLite Backup
Cron job harian (3 AM), `sqlite3 .backup` (online, no lock), gzip, keep 7 hari.

## Monorepo Structure

```
okuru.id/
├── backend/          # Goravel
├── admin/            # Vue + shadcn-vue
├── deploy/           # Dockerfile, docker-compose, Caddyfile
├── scripts/          # backup.sh, setup-vps.sh
├── docs/plans/       # design docs
├── .gitlab-ci.yml
└── Makefile
```

## Out of Scope (Phase 2)

- Payment integration detail (Sumopod API specifics) — ditunda, design terpisah
- Email delivery automation untuk digital product
- Analytics dashboard
- Multi-admin / role-based access
- Comments system untuk blog
