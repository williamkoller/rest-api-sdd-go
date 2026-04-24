# rest-api-sdd-go Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-24

## Active Technologies
- PostgreSQL 16 (primary), Redis 7 (optional cache) (HEAD)
- Go 1.26.2 + golangci-lint (external), govulncheck (`golang.org/x/vuln`), lefthook (002-makefile-fmt-lint-vulnerability-lefthook)

- Go 1.26.2 + Gin (HTTP), GORM + gorm/driver/postgres (ORM), Viper (config), go-redis/v9 (Redis cache), golang-jwt/jwt/v5 (auth), golang-migrate (DB migrations) (001-school-management-system)

## Project Structure

```text
cmd/api/main.go                              # Entry point
config/config.go                             # Viper config (ONLY place Viper is used)
internal/domain/entity/                      # Pure Go business entities (no external imports)
internal/domain/repository/                  # Repository interfaces
internal/application/usecase/                # Business logic (only domain/* imports)
internal/infrastructure/cache/               # Cache interface + memory/redis implementations
internal/infrastructure/database/            # GORM connection + pool setup
internal/infrastructure/repository/          # GORM implementations + model/ subdirectory
internal/transport/http/handler/             # Gin handlers
internal/transport/http/middleware/          # auth, tenant, logger, recovery
internal/transport/http/response/            # JSON envelope helpers
internal/transport/http/router.go            # Route registration
migrations/                                  # Plain SQL (NNN_description.up/down.sql)
docker/                                      # Dockerfile + docker-compose.yml
```

## Import Rules (enforced)

| Layer | Allowed imports |
|-------|----------------|
| `domain/entity` | stdlib only |
| `domain/repository` | `context`, `domain/entity` |
| `application/usecase` | `domain/*` |
| `infrastructure/*` | `domain/*`, external packages |
| `transport/http` | `application/usecase`, `gin`, `log/slog` |
| `config` | `viper`, stdlib |

- GORM types MUST NOT appear outside `internal/infrastructure/`
- Gin types MUST NOT appear outside `internal/transport/http/`
- `context.Context` is the first parameter of every I/O function
- All errors crossing a layer boundary MUST be wrapped: `fmt.Errorf("layer: %w", err)`

## Commands

```bash
# Run locally (infra via Docker)
docker compose -f docker/docker-compose.yml up postgres redis -d
go run ./cmd/api/...

# Run everything
docker compose -f docker/docker-compose.yml up --build

# Migrations
docker compose -f docker/docker-compose.yml run --rm migrate

# New migration
migrate create -ext sql -dir migrations -seq description_here
```

## Code Style

- No unit tests (per project constitution)
- Use `log/slog` for all structured logging
- Response envelope: `{"data": ..., "error": null, "meta": ...}` — use helpers in `transport/http/response/`
- Cache: always program to the `cache.Cache` interface, never to a concrete implementation
- Multi-tenancy: always scope DB queries by `school_id` extracted from JWT via context

## Recent Changes
- 002-makefile-fmt-lint-vulnerability-lefthook: Added Go 1.26.2 + golangci-lint (external), govulncheck (`golang.org/x/vuln`), lefthook
- HEAD: Added Go 1.26.2 + Gin (HTTP), GORM + gorm/driver/postgres (ORM), Viper (config), go-redis/v9 (Redis cache), golang-jwt/jwt/v5 (auth), golang-migrate (DB migrations)

- 001-school-management-system: Added Go 1.26.2 + Gin (HTTP), GORM + gorm/driver/postgres (ORM), Viper (config), go-redis/v9 (Redis cache), golang-jwt/jwt/v5 (auth), golang-migrate (DB migrations)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
