# Quickstart: School Management System API

**Feature**: `001-school-management-system`
**Date**: 2026-04-23

---

## Prerequisites

- Go 1.26.2+
- Docker + Docker Compose v2
- `make` (optional, for convenience)

---

## Repository Layout

```text
.
├── cmd/
│   └── api/
│       └── main.go                # Entry point
├── config/
│   └── config.go                  # Viper config loading
├── internal/
│   ├── domain/
│   │   ├── entity/                # Pure Go business entities
│   │   └── repository/            # Repository interfaces
│   ├── application/
│   │   └── usecase/               # Use cases per module
│   ├── infrastructure/
│   │   ├── cache/                 # Cache interface + memory/redis impls
│   │   ├── database/              # GORM setup + connection
│   │   └── repository/            # GORM repository implementations
│   │       └── model/             # GORM models (infra-only)
│   └── transport/
│       └── http/
│           ├── handler/           # Gin route handlers per module
│           ├── middleware/        # Auth, logger, tenant, recovery
│           └── router.go
├── migrations/                    # Plain SQL files (NNN_name.up.sql)
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
└── go.mod
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in:

```bash
# Application
APP_ENV=development          # development | production
APP_PORT=8080
APP_LOG_LEVEL=info           # debug | info | warn | error

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=escola_gestao
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Cache
CACHE_DRIVER=memory          # memory | redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_DEFAULT_TTL=300        # seconds

# JWT
JWT_SECRET_KEY=change-me-in-production
JWT_ACCESS_TOKEN_TTL=900     # 15 minutes in seconds
JWT_REFRESH_TOKEN_TTL=604800 # 7 days in seconds

# Migrations
RUN_MIGRATIONS=true          # auto-run on startup (dev only)
```

---

## Running with Docker Compose

```bash
# Start all services (postgres, redis, api)
docker compose -f docker/docker-compose.yml up

# First run: migrations run automatically via the migrate service
# API is available at http://localhost:8080

# Stop
docker compose -f docker/docker-compose.yml down

# Rebuild after code changes
docker compose -f docker/docker-compose.yml up --build
```

---

## Running Locally (without Docker)

```bash
# 1. Start infrastructure only
docker compose -f docker/docker-compose.yml up postgres redis -d

# 2. Install dependencies
go mod download

# 3. Run migrations
docker compose -f docker/docker-compose.yml run --rm migrate

# 4. Start the API
go run ./cmd/api/...
```

---

## Health Check

```bash
curl http://localhost:8080/health
# {"data":{"status":"ok","version":"1.0.0","db":"ok","cache":"ok"}}
```

---

## Running Migrations Manually

```bash
# Apply all pending migrations
docker compose -f docker/docker-compose.yml run --rm migrate

# Or directly with golang-migrate (if installed locally):
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/escola_gestao?sslmode=disable" up

# Roll back last migration
migrate -path ./migrations -database "..." down 1
```

---

## Creating a New Migration

```bash
# Creates migrations/NNN_description.up.sql and .down.sql
migrate create -ext sql -dir migrations -seq description_here
```

---

## Cache Strategy Switching

```bash
# In-memory (default for development)
CACHE_DRIVER=memory go run ./cmd/api/...

# Redis
CACHE_DRIVER=redis REDIS_ADDR=localhost:6379 go run ./cmd/api/...
```

---

## First API Call: Login

```bash
# 1. Create the first super-admin (via seed or direct DB insert in dev)
# See migrations/seed_superadmin.sql

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"changeme"}'

# 3. Use the access_token in subsequent requests
export TOKEN="<access_token from response>"

# 4. Create a school
curl -X POST http://localhost:8080/api/v1/schools \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Escola Exemplo","cnpj":"12345678000199","email":"contato@escola.com"}'
```

---

## Code Conventions

| Layer | Import rule |
|-------|-------------|
| `domain/entity` | No external imports |
| `domain/repository` | Only `context`, `domain/entity` |
| `application/usecase` | Only `domain/*` |
| `infrastructure/*` | `domain/*`, external packages |
| `transport/http` | Only `application/usecase`, `gin`, `log/slog` |
| `config` | Only `viper`, stdlib |

- GORM types (`*gorm.DB`, GORM model tags) MUST NOT appear outside `internal/infrastructure/`.
- Gin types (`*gin.Context`) MUST NOT appear outside `internal/transport/http/`.
- `context.Context` is the first parameter of every function touching I/O.
- All errors crossing a layer boundary MUST be wrapped: `fmt.Errorf("usecase: %w", err)`.

---

## Dev Tools Setup (branch `002-makefile-fmt-lint-vulnerability-lefthook`)

### One-time setup

```bash
# Install Go-based dev tools
make tools

# Install golangci-lint (external binary)
# macOS
brew install golangci-lint
# Linux / CI
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Activate git hooks
make hooks
```

### Daily usage

```bash
make fmt            # Format all Go code
make lint           # Run all linters
make vulnerability  # Scan for known vulnerabilities
```

### Commit message format

All commits are validated by the `commit-msg` lefthook:

```
type(scope)?: description
```

| Field | Required | Values |
|-------|----------|--------|
| type | yes | feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert |
| scope | no | lowercase alphanumeric + `/`, `_`, `-` (e.g. `auth`, `financial/invoice`) |
| description | yes | imperative, lowercase, no period |

```bash
git commit -m "feat(auth): add JWT refresh token endpoint"
git commit -m "fix: handle nil pointer in invoice generation"
git commit -m "chore(deps): bump golang.org/x/vuln"
```
