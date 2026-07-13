.PHONY: dev-frontend dev-backend dev build-local docker-up docker-down docker-logs deploy migrate seed backup build-frontend

# === LOCAL DEV ===
dev-frontend:
	cd frontend && bun run dev

dev-backend:
	@fuser -k 3000/tcp 2>/dev/null || true
	cd backend && go run .

dev:
	@make dev-backend & make dev-frontend & wait

build-frontend:
	cd frontend && bun run build && cp -r dist/* ../backend/public/

# === DOCKER ===
build-local:
	docker compose -f deploy/docker-compose.yml build

docker-up:
	docker compose -f deploy/docker-compose.yml up -d

docker-down:
	docker compose -f deploy/docker-compose.yml down

docker-logs:
	docker compose -f deploy/docker-compose.yml logs -f

# === DATABASE ===
migrate:
	docker compose -f deploy/docker-compose.yml exec app ./artisan migrate

seed:
	docker compose -f deploy/docker-compose.yml exec app ./artisan db:seed

# === DEPLOY ===
deploy:
	docker compose -f deploy/docker-compose.yml pull
	docker compose -f deploy/docker-compose.yml up -d --remove-orphans
	docker compose -f deploy/docker-compose.yml exec -T app ./artisan migrate --force

# === BACKUP ===
backup:
	docker compose -f deploy/docker-compose.yml exec app /opt/okuru/scripts/backup.sh
