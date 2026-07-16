.PHONY: dev-frontend dev-backend dev build docker-build docker-up docker-down docker-logs deploy migrate seed backup build-frontend

# === LOCAL DEV ===
dev-frontend:
	cd frontend && bun run dev

dev-backend:
	@fuser -k 3000/tcp 2>/dev/null || true
	cd backend && go run .

dev:
	@make dev-backend & make dev-frontend & wait

# Build frontend and copy output into backend/public/ so the Go server serves it.
build build-frontend:
	cd frontend && bun run build && cp -r dist/* ../backend/public/

# === DOCKER ===
docker-build:
	docker compose -f docker-compose.yml build

docker-up:
	docker compose -f docker-compose.yml up -d

docker-down:
	docker compose -f docker-compose.yml down

docker-logs:
	docker compose -f docker-compose.yml logs -f

# === DATABASE ===
migrate:
	docker compose -f docker-compose.yml exec app ./artisan migrate

seed:
	docker compose -f docker-compose.yml exec app ./artisan db:seed

# === DEPLOY ===
deploy:
	docker compose -f docker-compose.yml pull
	docker compose -f docker-compose.yml up -d --remove-orphans
	docker compose -f docker-compose.yml exec -T app ./artisan migrate --force

# === BACKUP ===
backup:
	docker compose -f docker-compose.yml exec app /opt/okuru/scripts/backup.sh
