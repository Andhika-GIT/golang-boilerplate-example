# Makefile for application commands
.PHONY: help

# ============================================
# DOCKER COMPOSE COMMANDS
# ============================================

.PHONY: dev-compose-build
dev-compose-build:
	@echo "🚀 Building development containers..."
	APP_ENV_FILE=.env APP_PORT=3005 docker compose build

.PHONY: dev-compose-up
dev-compose-up:
	@echo "🚀 Starting development containers..."
	APP_ENV_FILE=.env APP_PORT=3005 docker compose up -d

.PHONY: dev-compose-stop
dev-compose-stop:
	@echo "🛑 Stopping development containers..."
	docker compose stop

.PHONY: dev-compose-down
dev-compose-down:
	@echo "🛑 Stopping and removing development containers..."
	docker compose down

.PHONY: dev-logs
dev-logs:
ifdef c
	@echo "📋 Showing logs for container: $(c)"
	@docker compose logs -f $(c)
else
	@echo "📋 Showing all container logs..."
	@docker compose logs -f
endif

# ============================================
# DATABASE MIGRATION COMMANDS
# ============================================

.PHONY: migrate-up
migrate-up:
	@echo "🔄 Running database migrations..."
	@docker exec -it app-c-dev /app/tmp/main migrate up

.PHONY: migrate-down
migrate-down:
	@echo "⏪ Rolling back migrations..."
	@docker exec -it app-c-dev /app/tmp/main migrate down

.PHONY: migrate-status
migrate-status:
	@echo "📊 Checking migration status..."
	@docker exec -it app-c-dev /app/tmp/main migrate status

.PHONY: migrate-create
migrate-create:
	@echo "📝 Creating new migration: $(name)"
	@docker exec -it app-c-dev /app/tmp/main migrate create $(name)

# ============================================
# DATABASE SEEDER COMMANDS (DOCKER)
# ============================================

.PHONY: seed
seed:
	@echo "🌱 Running database seeders..."
	@docker exec -it app-c-dev /app/tmp/main seed

.PHONY: seed-specific
seed-specific:
	@echo "🌱 Running specific seeder..."
	@docker exec -it app-c-dev /app/tmp/main seed --name $(name)

.PHONY: seed-list
seed-list:
	@echo "📋 Listing available seeders..."
	@docker exec -it app-c-dev /app/tmp/main seed --list

# ============================================
# HELP
# ============================================

.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make dev-compose-build    - Build development containers"
	@echo "  make dev-compose-up       - Start development containers"
	@echo "  make dev-compose-stop     - Stop development containers"
	@echo "  make dev-compose-down     - Stop and remove development containers"
	@echo "  make dev-logs c=<service> - Show container logs (specific or all)"
	@echo ""
	@echo "Database Migrations:"
	@echo "  make migrate-up           - Run pending migrations"
	@echo "  make migrate-down         - Rollback migrations"
	@echo "  make migrate-status       - Show migration status"
	@echo "  make migrate-create name=<name> - Create new migration"
	@echo ""
	@echo "Database Seeders:"
	@echo "  make seed                 - Run all seeders"
	@echo "  make seed-specific name=<name> - Run specific seeder"
	@echo "  make seed-list            - List available seeders"