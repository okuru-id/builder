# Okuru

Personal platform for [okuru.id](https://okuru.id) — blog, portfolio, storefront, and
digital product sales. Single Go binary serves the public site, storefront, and the
JSON API consumed by the Vue admin SPA.

## Tech Stack

| Layer     | Technology                                            |
|-----------|-------------------------------------------------------|
| Backend   | [Goravel](https://goravel.dev) (Go 1.24), Gin        |
| Database  | SQLite (default, via `ncruces/go-sqlite3`) or Postgres |
| Admin SPA | Vue 3 + shadcn-vue (reka-ui) + Tailwind v4 + Vite    |
| Storefront| Go `html/template` + Tailwind (CDN) + Alpine.js       |
| Proxy/TLS | Caddy                                                 |
| CI/CD     | GitLab CI                                             |
| Deploys   | Docker Compose on a VPS                               |

## Repository Layout

```
.
├── admin/        # Vue 3 admin SPA (shadcn-vue)
├── backend/      # Goravel app: API, storefront, static admin host
├── deploy/       # Dockerfile, compose files, Caddyfile
├── docs/         # Plans and design notes
└── scripts/      # Operational scripts (DB backup, etc.)
```

## Local Development

Requirements: Go 1.24+, Node 20+, pnpm 10+.

```bash
# Backend (http://localhost:3000)
cd backend
cp .env.example .env   # if absent, defaults work for SQLite
go run .               # serves API + storefront templates

# Admin SPA (http://localhost:5173 -> proxies /admin/api to :3000)
cd admin
pnpm install
pnpm run dev
```

### Database

```bash
cd backend
./artisan migrate          # apply migrations
./artisan db:seed          # optional seeders
```

### Handy Make targets

```bash
make dev          # run backend (air) + admin concurrently
make build-local  # build the production Docker image
make docker-up    # start the full stack via compose
make migrate      # run migrations inside the container
make backup       # take a gzipped SQLite backup inside the container
```

## Docker

Local stack (builds the image from source):

```bash
make build-local && make docker-up
# API/admin: http://localhost:8080
# Storefront: http://localhost:8080/shop
# Health:     http://localhost:8080/health
```

Production stack pulls the GitLab registry image and fronts it with Caddy
(TLS for `api.okuru.id`, `shop.okuru.id`). See `deploy/docker-compose.prod.yml`.

## Project Conventions

- **Backend routes**: public API under `api/v1`, admin API under `admin/api`,
  storefront + health in `routes/web.go`.
- **Views**: Go templates live in `backend/resources/views/`, rendered via
  `ctx.Response().View().Make("path.tmpl", data)`.
- **Admin SPA**: built to `admin/dist`, copied to `backend/public/admin/` in the
  Docker image, served under `/admin/`.
- **Payments**: Sumopod integration is planned (Phase 2); checkout currently
  returns 501.

See [`AGENTS.md`](./AGENTS.md) for build/test commands and code-style guidance.
