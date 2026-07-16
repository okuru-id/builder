# AGENTS.md

Guidance for AI coding agents working in this repository.

## Build & Run Commands

### Backend (`backend/`)
```bash
go build ./...                  # compile-check
go run .                        # start server on :3000
go vet ./...                    # static checks
./artisan migrate               # apply DB migrations
./artisan db:seed               # run seeders
CGO_ENABLED=0 go build -o okuru .   # reproducible prod binary (ncruces sqlite = pure Go)
```

No CGO is required: SQLite uses `ncruces/go-sqlite3` (WASM).

### Frontend (`frontend/`)
```bash
bun install
bun run build      # vue-tsc typecheck + vite build -> frontend/dist
bun run dev        # dev server (landing at /, admin at /admin/), proxies /admin/api -> localhost:3000
```

### Docker
```bash
make docker-build    # build Docker image
make docker-up      # run the local compose stack
```

## Tech Stack
- **Goravel** (Gin driver) on Go 1.24. App entry: `backend/main.go` -> `bootstrap/app.go`.
- **ORM**: Goravel ORM (GORM-backed). Local facades in `backend/app/facades/`.
- **DB**: SQLite (default) or Postgres, configured in `backend/config/database.go`.
- **Admin**: Vue 3 + shadcn-vue (reka-ui) + Tailwind v4 + Vite. Package manager: bun. SPA mounted at `/admin/`.
- **Storefront**: server-rendered Go templates + Tailwind CDN + Alpine.js.

## Code Style

### Go
- Follow existing controller patterns in `backend/app/http/controllers/`.
- Controllers: struct + `New<Name>Controller()` constructor, methods take
  `http.Context` and return `http.Response`.
- ORM queries: `facades.Orm().Query().Where(...).First(&model)`. Note:
  `First()` may return a zero-value struct without error for missing rows — guard
  with `model.ID == 0` when a 404 is expected.
- JSON responses: `ctx.Response().Success().Json(http.Json{...})` or
  `ctx.Response().Json(http.StatusX, http.Json{...})`.
- Views: `ctx.Response().View().Make("shop/index.tmpl", data)`. Each template file
  must `{{ define "path.tmpl" }}...{{ end }}`.
- Public API under prefix `api/v1`; admin API under `admin/api` (JWT + TOTP middleware).
- Add comments on exported handlers (one line, English OK).

### Frontend (admin)
- TypeScript, `<script setup lang="ts">` SFCs.
- Components: shadcn-vue primitives from `frontend/src/components/ui/`.
- API calls via `axios`; routes in `frontend/src/router`.

## Routing Map
- `GET /health` -> health check
- `GET /shop`, `GET /product/{slug}`, `POST /checkout/{slug}` -> storefront
- `GET /admin/*` -> admin SPA static files
- `api/v1/*` -> public content API
- `admin/api/*` -> authenticated admin CRUD

## Migration & Models
- Migrations in `backend/database/migrations/` (timestamped). Run via `./artisan migrate`.
- Models in `backend/app/models/` embed `orm.Model` (provides `ID`, timestamps).
- Register new migrations in `backend/bootstrap/migrations.go`, seeders in
  `backend/bootstrap/seeders.go`, providers in `backend/bootstrap/providers.go`.

## Deployment
- Image: `Dockerfile` (3-stage: admin build -> Go build -> alpine runtime).
- Runtime serves the single `okuru` binary; `./artisan` wrapper calls
  `./okuru artisan` (no Go toolchain needed at runtime).
- CI: `.gitlab-ci.yml` -> lint, build/push image, deploy to VPS via SSH.
- TLS termination and virtual hosts handled by Caddy (`infra/Caddyfile`).

## Things to Avoid
- Do not add CGO dependencies.
- Do not introduce a separate frontend framework for the storefront (keep templates).
- Do not commit `backend/.env`, `admin/dist/`, or `backend/public/admin/`.
