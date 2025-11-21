# Makefile placeholder
.PHONY: help


.PHONY: dev-compose-build
dev-compose-build:
	@echo "🚀 Starting Compose Build For Development..."
	APP_ENV_FILE=.env APP_PORT=3005 docker compose build

.PHONY: dev-compose-up
dev-compose-up:
	@echo "🚀 Starting Compose Up For Development..."
	APP_ENV_FILE=.env APP_PORT=3005 docker compose up -d

.PHONY: dev-compose-Stop
dev-compose-stop:
	@echo "🛑 Stopping Development Containers..."
	docker compose stop

.PHONY: dev-compose-down
dev-compose-down:
	@echo "🛑 Stopping And Delete Development Containers..."
	docker compose down

dev-logs:
	@echo "📋 Showing Container Logs..."
	docker compose logs -f

help:
	@echo "Available commands:"
	@echo "  make run env=.env.dev port=3005    - Custom env and port"
	@echo "  make dev                          - Development (.env.dev)"
	@echo "  make development                  - Development (.env.development)"
	@echo "  make staging                      - Staging (.env.staging)"
	@echo "  make prod                         - Production (.env.production)"
	@echo "  make down                         - Stop and delete containers"
	@echo "  make clean                        - Stop and remove volumes"
