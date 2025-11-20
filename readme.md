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
│  └── migrate.go                         │
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
│  ├── redis.go                           │
│                        │
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
│       └── migrate.go             # Database migration commands
│
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loader
│   │
│   ├── connections/               # Low-level connections (pure functions)
│   │   ├── database.go            # Database connection functions
│   │   ├── redis.go               # Redis connection
│   │  
│   │
│   ├── infrastructure/            # High-level orchestration
│   │   ├── connections.go         # Connection orchestrator
│   │   └── di.go                  # Dependency injection
│   │
│   ├── database/
│   │   └── migrations/            # SQL migration files
│   │       ├── 000001_xxx.up.sql
│   │       └── 000001_xxx.down.sql
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
docker exec be-c-dev /app/tmp/main migrate up
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

---

## 🔧 Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **golang-migrate CLI** (installed in Docker, or locally for native development)

### Install golang-migrate locally (optional):

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Or via Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
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
# Start all services (PostgreSQL, Redis, App)
docker-compose up -d

# Check logs
docker-compose logs -f app
```

### 4. Run Migrations

```bash
# Inside Docker container
docker exec be-c-dev /app/tmp/main migrate up

# Check migration status
docker exec be-c-dev /app/tmp/main migrate status
```


## 💻 Development

### Docker Development (Recommended)

```bash
# Start services
docker-compose up -d

# View logs (with hot-reload via Air)
docker-compose logs -f app

# Execute commands in container
docker exec be-c-dev /app/tmp/main [command]

# Stop services
docker-compose down
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

### Migration Commands (Inside Docker)

```bash
# Create new migration
docker exec be-c-dev /app/tmp/main migrate create create_users_table

# Run all pending migrations
docker exec be-c-dev /app/tmp/main migrate up

# Rollback 1 migration
docker exec be-c-dev /app/tmp/main migrate down

# Rollback N migrations
docker exec be-c-dev /app/tmp/main migrate down 3

# Check migration status
docker exec be-c-dev /app/tmp/main migrate status
```

### Migration Workflow

#### 1. Create Migration

```bash
docker exec be-c-dev /app/tmp/main migrate create create_users_table
```

This creates two files:
- `internal/database/migrations/000001_create_users_table.up.sql`
- `internal/database/migrations/000001_create_users_table.down.sql`

#### 2. Write Migration SQL

**000001_create_users_table.up.sql:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**000001_create_users_table.down.sql:**
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

#### 3. Run Migration

```bash
docker exec be-c-dev /app/tmp/main migrate up
```

#### 4. Test Rollback (Optional)

```bash
docker exec be-c-dev /app/tmp/main migrate down
docker exec be-c-dev /app/tmp/main migrate up
```

---

## 📖 Commands Reference

### Server Commands

```bash
# Start server (default port 8080)
docker exec be-c-dev /app/tmp/main server

# Start with custom port
docker exec be-c-dev /app/tmp/main server --port 9000
```

### Migration Commands

```bash
# Create new migration
docker exec be-c-dev /app/tmp/main migrate create [name]

# Run all pending migrations
docker exec be-c-dev /app/tmp/main migrate up

# Rollback migrations
docker exec be-c-dev /app/tmp/main migrate down [steps]

# Show current migration version
docker exec be-c-dev /app/tmp/main migrate status
```

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d db

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app
docker-compose logs -f db

# Rebuild containers
docker-compose up -d --build

# Access container shell
docker exec -it be-c-dev sh

# Access PostgreSQL
docker exec -it db-c-dev psql -U postgres -d go_stream_dev

# Access Redis
docker exec -it redis-c-dev redis-cli
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
docker-compose down -v
docker-compose up -d
docker exec be-c-dev /app/tmp/main migrate up
```

### Add New Feature with Migration

```bash
# 1. Create migration
docker exec be-c-dev /app/tmp/main migrate create add_feature_x

# 2. Edit SQL files
# internal/database/migrations/00000X_add_feature_x.up.sql
# internal/database/migrations/00000X_add_feature_x.down.sql

# 3. Run migration
docker exec be-c-dev /app/tmp/main migrate up

# 4. Code changes will auto-reload via Air
# Edit your Go files...

# 5. Test
curl http://localhost:3005/api/v1/...
```

### Database Reset

```bash
# Rollback all migrations
docker exec be-c-dev /app/tmp/main migrate down 999

# Or recreate container
docker-compose down -v db
docker-compose up -d db

# Wait for DB ready, then migrate
sleep 3
docker exec be-c-dev /app/tmp/main migrate up
```

---

## 🔍 Troubleshooting

### App can't connect to database

```bash
# Check if DB is running
docker ps | grep db-c-dev

# Check DB logs
docker-compose logs db

# Verify DB connection from app container
docker exec be-c-dev ping host.docker.internal
```

### Migration fails

```bash
# Check migration status
docker exec be-c-dev /app/tmp/main migrate status

# Check if migrate CLI is available
docker exec be-c-dev migrate -version

# View migration files
docker exec be-c-dev ls -la /app/internal/database/migrations/
```

### Hot-reload not working

```bash
# Check Air logs
docker-compose logs -f app

# Restart app container
docker-compose restart app

# Check .air.toml configuration
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

---

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Create migration if needed
4. Test locally
5. Submit pull request
