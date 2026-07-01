.PHONY: dev-admin dev-backend dev build-local docker-up docker-down docker-logs deploy migrate seed backup

# === LOCAL DEV ===
dev-admin:
	cd admin && pnpm run dev

dev-backend:
	cd backend && air

dev:
	@make dev-backend & make dev-admin & wait

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
