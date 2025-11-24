# Simple Go Microservice BoilerPlate

A production-ready Go microservice with clean architecture, built with Gin, Cobra CLI, and golang-migrate.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Application Flow](#application-flow)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Development](#development)
- [Database Migrations](#database-migrations)
- [Database Seeders](#database-seeders)
- [Commands Reference](#commands-reference)
- [Configuration](#configuration)

---

## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│           cmd/main.go                   │  Entry Point
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│      cmd/commands/ (Cobra CLI)          │  CLI Layer
│  ├── root.go                            │
│  ├── server.go                          │
│  ├── migrate.go                         │
│  └── seeder.go                          │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│      internal/config/                   │  Configuration
│  └── config.go                          │  (Load from env)
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│   internal/infrastructure/              │  Infrastructure
│  ├── connections.go                     │  (Orchestrator)
│  └── di.go                              │  (Dependency Injection)
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│   internal/connections/                 │  Low-Level Connections
│  ├── database.go                        │  (Pure functions)
│  └── redis.go                           │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│   Business Logic Layers                 │
│  ├── repositories/                      │  Data Access
│  ├── services/                          │  Business Logic
│  ├── controllers/                       │  HTTP Handlers
│  └── routes/                            │  Route Registry
└─────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
user-service/
├── cmd/
│   ├── main.go                    # Entry point
│   └── commands/                  # Cobra CLI commands
│       ├── root.go                # Root command & registration
│       ├── server.go              # HTTP server command
│       ├── migrate.go             # Database migration commands
│       └── seeder.go              # Database seeder commands
│
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loader
│   │
│   ├── connections/               # Low-level connections (pure functions)
│   │   ├── database.go            # Database connection functions
│   │   └── redis.go               # Redis connection
│   │
│   ├── infrastructure/            # High-level orchestration
│   │   ├── connections.go         # Connection orchestrator
│   │   └── di.go                  # Dependency injection
│   │
│   ├── database/
│   │   ├── migrations/            # SQL migration files
│   │   │   ├── 000001_xxx.up.sql
│   │   │   └── 000001_xxx.down.sql
│   │   └── seeders/               # Database seeders
│   │       ├── seeder.go          # Seeder registry
│   │       ├── roles.seed.go
│   │       └── users.seed.go
│   │
│   ├── models/                    # Data models
│   │   ├── user.go
│   │   └── role.go
│   │
│   ├── repositories/              # Data access layer
│   │   ├── interfaces.go
│   │   └── user_repository.go
│   │
│   ├── services/                  # Business logic layer
│   │   ├── interfaces.go
│   │   └── user_service.go
│   │
│   ├── controllers/               # HTTP handlers
│   │   ├── registry.go
│   │   └── user/
│   │       └── user_controller.go
│   │
│   ├── routes/                    # Route definitions
│   │   ├── registry.go
│   │   └── user/
│   │       └── user_routes.go
│   │
│   └── shared/                    # Shared utilities
│       ├── constants/
│       ├── middleware/
│       └── commons/
│
├── .env                           # Environment variables
├── .air.toml                      # Air hot-reload config
├── docker-compose.yml             # Docker services
├── dev.Dockerfile                 # Development Dockerfile
├── Makefile                       # Development commands
└── go.mod
```

---

## 🔄 Application Flow

### 1. **Startup Flow (Server)**

```
cmd/main.go
    ↓
commands.Execute()
    ↓
commands/server.go → runServer()
    ↓
config.Load() ──────────────→ Load environment variables
    ↓
infrastructure.InitializeConnections(cfg)
    ↓
    ├─→ Read cfg.Database.Type (postgres/mysql)
    ├─→ Call connections.NewPostgresConnection() or NewMySQLConnection()
    ├─→ Configure connection pool
    ├─→ Initialize Redis (if enabled)
    └─→ Return Connections{DB, Redis}
    ↓
infrastructure.InitializeDependencies(connections, cfg)
    ↓
    ├─→ Create Repositories (inject DB connection)
    ├─→ Create Services (inject repositories)
    ├─→ Create Controllers (inject services)
    └─→ Return Dependencies{Services, Controllers}
    ↓
routes.NewRouteRegistry(router, deps.ControllerRegistry)
    ↓
    └─→ Register all HTTP routes
    ↓
router.Run() ──────────────→ Start HTTP server
```

### 2. **Request Flow**

```
HTTP Request
    ↓
Gin Router (routes/)
    ↓
Controller (controllers/)
    ↓
Service (services/) ─────→ Business Logic
    ↓
Repository (repositories/) ─→ Database Query
    ↓
Database (PostgreSQL/MySQL)
    ↓
← Response flows back ←
```

### 3. **Migration Flow**

```
make migrate-up
    ↓
commands/migrate.go → runMigrateUp()
    ↓
getDatabaseURL() ──────────→ Build connection URL from config
    ↓
exec.Command("migrate", ...) ─→ Execute golang-migrate CLI
    ↓
    └─→ Read SQL files from internal/database/migrations/
    ↓
Apply migrations to database
```

### 4. **Seeder Flow**

```
make seed
    ↓
commands/seeder.go → runSeeder()
    ↓
config.Load() ──────────────→ Load configuration
    ↓
connections.NewDatabase() ──→ Initialize database connection
    ↓
seeders.InitSeeders(db) ────→ Register all seeders
    ↓
registry.RunAll() ──────────→ Execute seeders in order
    ↓
Seeders insert data to database
```

---

## 🔧 Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **Make** (usually pre-installed on Linux/macOS)
- **golang-migrate CLI** (installed in Docker)

### Install Make (if needed):

```bash
# macOS (usually pre-installed)
xcode-select --install

# Ubuntu/Debian
sudo apt-get install build-essential

# Windows
choco install make
```

---

## 🚀 Getting Started

### 1. Clone Repository

```bash
git clone <repository-url>
cd user-service
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Start Services

```bash
# Build and start all services (PostgreSQL, Redis, App)
make dev-compose-build
make dev-compose-up

# Check logs
make dev-logs

# Check specific service logs
make dev-logs c=app
```

### 4. Run Migrations

```bash
# Run all pending migrations
make migrate-up

# Check migration status
make migrate-status
```

### 5. Seed Database (Optional)

```bash
# Run all seeders
make seed

# Or run specific seeder
make seed-specific name=RoleSeeder
```

---

## 💻 Development

### Using Makefile (Recommended)

All common tasks can be executed via Makefile commands:

```bash
# View all available commands
make help

# Start development environment
make dev-compose-up

# View logs (hot-reload via Air)
make dev-logs

# View specific service logs
make dev-logs c=app
make dev-logs c=db

# Stop services
make dev-compose-stop

# Stop and remove containers
make dev-compose-down
```

### Local Development (Without Docker)

```bash
# Install dependencies
go mod download

# Run server
go run cmd/main.go server

# Or with Air (hot-reload)
air -c .air.toml
```

---

## 🗄️ Database Migrations

### Migration Commands (Using Makefile)

```bash
# Create new migration
make migrate-create name=create_users_table

# Run all pending migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status
```

### Migration Workflow

#### 1. Create Migration

```bash
make migrate-create name=create_users_table
```

This creates two files:
- `internal/database/migrations/000001_create_users_table.up.sql`
- `internal/database/migrations/000001_create_users_table.down.sql`

#### 2. Write Migration SQL

**000001_create_users_table.up.sql:**
```sql
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(250) UNIQUE,
    phone VARCHAR(23) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_users_email_verified ON users(email_verified);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);
);

CREATE INDEX idx_users_email ON users(email);
```

**000001_create_users_table.down.sql:**
```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email_verified;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_created_at;

-- Drop table
DROP TABLE IF EXISTS users;
```

#### 3. Run Migration

```bash
make migrate-up
```

#### 4. Test Rollback (Optional)

```bash
make migrate-down
make migrate-up
```

---

## 🌱 Database Seeders

Seeders are used to populate the database with initial or test data.

### Seeder Commands (Using Makefile)

```bash
# List all available seeders
make seed-list

# Run all seeders
make seed

# Run specific seeder
make seed-specific name=RoleSeeder
make seed-specific name=UserSeeder
```

### Creating New Seeder

1. Create a new file in `internal/database/seeders/`:

```go
// internal/database/seeders/categories.seed.go
package seeders

import (
    "gorm.io/gorm"
    "your-project/internal/models"
)

type CategorySeeder struct{
    db *gorm.DB
}

func (s *CategorySeeder) GetName() string {
    return "CategorySeeder"
}

func (s *CategorySeeder) Seed() error {
    categories := []models.Category{
        {Name: "Technology", Slug: "technology"},
        {Name: "Business", Slug: "business"},
    }

    for _, category := range categories {
        var existing models.Category
        if err := s.db.Where("slug = ?", category.Slug).First(&existing).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                if err := s.db.Create(&category).Error; err != nil {
                    return err
                }
            } else {
                return err
            }
        }
    }
    
    return nil
}
```

2. Register the seeder in `internal/database/seeders/seeder.go`:

```go
func InitSeeders(db *gorm.DB) *SeederRegistry {
    registry := NewSeederRegistry(db)
    
    registry.Register(&RoleSeeder{db: db})
    registry.Register(&UserSeeder{db: db})
    registry.Register(&CategorySeeder{db: db}) // Add this line
    
    return registry
}
```

3. Run the seeder:

```bash
make seed
# or
make seed-specific name=CategorySeeder
```

---

## 📖 Commands Reference

### Docker Compose Commands

```bash
# Build containers
make dev-compose-build

# Start containers
make dev-compose-up

# Stop containers (keeps data)
make dev-compose-stop

# Stop and remove containers (removes data)
make dev-compose-down

# View logs
make dev-logs                    # All services
make dev-logs c=app             # App only
make dev-logs c=db              # Database only
```

### Migration Commands

```bash
# Create new migration
make migrate-create name=<migration_name>

# Run all pending migrations
make migrate-up

# Rollback migrations
make migrate-down

# Show current migration version
make migrate-status
```

### Seeder Commands

```bash
# List all available seeders
make seed-list

# Run all seeders
make seed

# Run specific seeder
make seed-specific name=<SeederName>
```

### Manual Docker Commands (Alternative)

```bash
# Access container shell
docker exec -it app-c-dev sh

# Access PostgreSQL
docker exec -it db-c-dev psql -U postgres -d go_stream_dev

# Access Redis
docker exec -it redis-c-dev redis-cli

# Direct command execution
docker exec -it app-c-dev /app/tmp/main server
docker exec -it app-c-dev /app/tmp/main migrate status
docker exec -it app-c-dev /app/tmp/main seed
```

---

## ⚙️ Configuration

Configuration is loaded from environment variables via `.env` file.

### Key Configuration Variables

```bash
# Application
APP_NAME=user-service
APP_VERSION=1.0.0
APP_ENV=development
APP_DEBUG=true

# Server
SERVER_PORT=3005

# Database
DB_TYPE=postgres              # postgres, mysql
DB_HOST=host.docker.internal  # Use this for Docker → Host DB
DB_PORT=5440                  # Host port mapped to container
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_stream_dev
DB_SSL_MODE=disable

# Redis
REDIS_ENABLED=true
REDIS_HOST=redis
REDIS_PORT=6379

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

### Database Connection in Docker

When running app in Docker and database on host machine:

```bash
# Use special DNS name
DB_HOST=host.docker.internal
DB_PORT=5440  # Host port (not container port 5432!)
```

When database is also in Docker:

```bash
# Use service name from docker-compose.yml
DB_HOST=db
DB_PORT=5432  # Container port
```

---

## 🏃 Common Workflows

### Fresh Start

```bash
# Clean everything and start fresh
make dev-compose-down
make dev-compose-build
make dev-compose-up
make migrate-up
make seed
```

### Add New Feature with Migration

```bash
# 1. Create migration
make migrate-create name=add_feature_x

# 2. Edit SQL files
# internal/database/migrations/00000X_add_feature_x.up.sql
# internal/database/migrations/00000X_add_feature_x.down.sql

# 3. Run migration
make migrate-up

# 4. Code changes will auto-reload via Air
# Edit your Go files...

# 5. Test
curl http://localhost:3005/api/v1/...
```

### Database Reset

```bash
# Rollback all migrations
make migrate-down

# Or recreate container with fresh data
make dev-compose-down
make dev-compose-up

# Wait for DB ready, then migrate and seed
sleep 3
make migrate-up
make seed
```

### Development Cycle

```bash
# Start development
make dev-compose-up
make dev-logs

# Make code changes (auto-reload via Air)
# ...

# Check logs
make dev-logs c=app

# Restart if needed
make dev-compose-stop
make dev-compose-up
```

---

## 🔍 Troubleshooting

### App can't connect to database

```bash
# Check if DB is running
docker ps | grep db-c-dev

# Check DB logs
make dev-logs c=db

# Verify DB connection from app container
docker exec app-c-dev ping host.docker.internal
```

### Migration fails

```bash
# Check migration status
make migrate-status

# Check if migrate CLI is available
docker exec app-c-dev migrate -version

# View migration files
docker exec app-c-dev ls -la /app/internal/database/migrations/
```

### Seeder fails

```bash
# List available seeders
make seed-list

# Check database connection
make dev-logs c=db

# Run specific seeder for debugging
make seed-specific name=RoleSeeder
```

### Hot-reload not working

```bash
# Check Air logs
make dev-logs c=app

# Restart app container
make dev-compose-stop
make dev-compose-up

# Check .air.toml configuration
```

### Make command not found

```bash
# Install make
# macOS
xcode-select --install

# Ubuntu/Debian
sudo apt-get install build-essential

# Windows
choco install make
```

---

## 📚 Additional Resources

### Architecture Patterns

- **Clean Architecture**: Business logic independent of frameworks
- **Dependency Injection**: Dependencies injected from outside
- **Repository Pattern**: Abstract data access layer
- **Service Layer**: Encapsulate business logic

### Tools Used

- **Gin**: HTTP web framework
- **Cobra**: CLI framework
- **golang-migrate**: Database migrations
- **Air**: Hot-reload for development
- **Docker**: Containerization
- **Make**: Task automation

### Project Benefits

- ✅ Hot-reload development with Air
- ✅ Easy command execution via Makefile
- ✅ Clean architecture with DI
- ✅ Database migrations and seeders
- ✅ Docker containerization
- ✅ Ready for production deployment

---

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Create migration if needed
4. Create/update seeders if needed
5. Test locally with `make dev-compose-up`
6. Submit pull request

---

## 📝 Quick Reference Card

```bash
# Essential Commands
make help                     # Show all commands
make dev-compose-up          # Start development
make dev-logs                # View logs
make migrate-up              # Run migrations
make seed                    # Populate database
make dev-compose-down        # Stop everything

# Development Workflow
make dev-compose-build       # Build containers
make migrate-create name=X   # Create migration
make seed-specific name=X    # Run specific seeder
make migrate-status          # Check migrations
make seed-list               # List seeders
```
