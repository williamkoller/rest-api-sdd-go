# Research: School Management System

**Feature**: `001-school-management-system`
**Date**: 2026-04-23
**Status**: Complete — all NEEDS CLARIFICATION resolved

---

## Resolved Clarifications

### Q1: Multi-tenant SaaS vs Single Organization

**Decision**: Multi-tenant SaaS — multiple independent schools share the platform.

**Rationale**: The spec's data isolation requirement (SC-008) and the School → Unit hierarchy
strongly imply separate organizations. Designing for multi-tenancy from the start is cheaper
than retrofitting it later. Row-level tenant isolation via `school_id` on all tables is the
chosen strategy (simpler than schema-per-tenant for this scale).

**Alternatives considered**:
- Single-organization: simpler, but limits commercial scope and complicates future expansion.
- Schema-per-tenant: stronger isolation, but adds operational complexity (migrations, routing).

---

### Q2: "Fila de Escola" Interpretation

**Decision**: Enrollment waitlist for prospective students (Option A).

**Rationale**: Most common interpretation in Brazilian school management context. No real-time
infrastructure (WebSockets, SSE) is needed, keeping the system simple.

**Alternatives considered**:
- Real-time pickup queue: would require WebSocket or SSE, significantly expanding scope and
  violating the Simplicity principle.

---

## Technology Stack Research

### Go 1.26.2

**Decision**: Use Go 1.26.2 as specified.

**Key features available** (cumulative from 1.21+):
- `log/slog`: structured logging (Principle V)
- `slices`, `maps` packages: idiomatic collection operations
- `errors.Join`: multi-error composition
- Range-over-func iterators (1.22+)
- Generic type inference improvements
- `math/rand/v2` for modern random generation

**Rationale**: Latest stable release with all modern language features.

---

### Gin (HTTP Framework)

**Decision**: `github.com/gin-gonic/gin` latest stable.

**Architecture boundary**: Gin is strictly confined to `internal/transport/http/`. No Gin type
(`*gin.Context`, `gin.HandlerFunc`) may leak into `internal/application/` or `internal/domain/`.

**Key patterns**:
- Route groups for versioning: `/api/v1/`
- Middleware chain for: auth (JWT), tenant scoping, structured logging, recovery
- Request binding via `ShouldBindJSON` + manual validation (no heavy validator packages)
- Response helpers standardize `{"data": ..., "error": null}` envelope

**Alternatives considered**: `net/http` (stdlib) — simpler but adds significant boilerplate
for routing, middleware, and binding for a project with 50+ endpoints.

---

### GORM + PostgreSQL

**Decision**: `gorm.io/gorm` with `gorm.io/driver/postgres`.

**Architecture notes** (justification for constitution violation):
- GORM models are defined ONLY in `internal/infrastructure/repository/model/` with GORM struct
  tags. They never appear in domain or application layers.
- Domain entities in `internal/domain/entity/` are plain Go structs with no GORM dependency.
- Mapper functions in each repository implementation translate between GORM models ↔ domain
  entities.
- This ensures the domain layer has zero dependency on GORM; switching ORM is possible by
  replacing infrastructure implementations only.

**Schema management**: GORM AutoMigrate is ONLY used in development (`APP_ENV=development`).
Production uses plain SQL migration files in `migrations/` applied by `golang-migrate`.

**Connection pooling**: Configured via `sql.DB` settings exposed by `gorm.DB.DB()`:
- `MaxOpenConns`: 25
- `MaxIdleConns`: 5
- `ConnMaxLifetime`: 5 minutes

**Alternatives considered**: `sqlc` + plain SQL — eliminates ORM but adds ~2000 lines of
generated/manual mapping code for 20+ entities. Rejected because no unit tests exist to
catch SQL query regressions.

---

### Viper (Configuration)

**Decision**: `github.com/spf13/viper` for configuration loading.

**Architecture notes** (justification for constitution violation):
- Viper is used ONLY in `config/config.go`. No other package imports Viper.
- A single `Config` struct is populated at startup and injected as a value throughout the app.
- All config values come from environment variables (`viper.AutomaticEnv()`). No config files
  in production.
- Viper's value-add: env binding to struct fields, type coercion, and defaults — replacing
  ~150 lines of `os.Getenv` + `strconv` boilerplate for 8+ configuration domains.

**Configuration domains**: DB, Redis, HTTP server, JWT, SMTP (notifications), cache driver.

**Alternatives considered**: Raw `os.Getenv` — simpler but verbose; no type safety or defaults
without manual boilerplate. Acceptable for 2-3 env vars, not for 30+.

---

### Redis Cache (with Strategy Pattern)

**Decision**: `github.com/redis/go-redis/v9` for Redis; `github.com/patrickmn/go-cache`
for in-memory (or `sync.Map` for zero-dependency option).

**Cache interface** (defined in `internal/infrastructure/cache/`):
```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    DeletePattern(ctx context.Context, pattern string) error
}
```

**Strategy selection**: `CACHE_DRIVER=memory|redis` environment variable. Factory function in
`config/` returns the appropriate implementation. Application code depends only on the
interface — no conditional imports in use-case or handler code.

**Cache invalidation strategy**: Write-through on mutation; cache-aside on read. TTL per
resource type (configurable). Pattern-based invalidation for collection caches.

**In-memory trade-offs**: No persistence, no shared state across instances. Acceptable for
single-instance dev and small deployments. Redis required for multi-instance production.

---

### Docker / Containers

**Decision**: Multi-stage Dockerfile. `docker-compose.yml` for local development.

**Services** (docker-compose):
- `api`: Go application
- `postgres`: PostgreSQL 16
- `redis`: Redis 7 (optional for dev, replaceable by memory cache)
- `migrate`: `golang-migrate` runner (runs on startup, then exits)

**Multi-stage Dockerfile**:
1. `builder`: `golang:1.26-alpine` — compiles binary
2. `final`: `alpine:latest` — minimal runtime image with CA certs only

---

### Authentication

**Decision**: JWT (RS256) with refresh token rotation.

**Flow**: Login → access token (15 min) + refresh token (7 days). Refresh endpoint issues
new pair. Logout invalidates refresh token in Redis (or memory cache).

**JWT claims**: `user_id`, `school_id`, `role`, `exp`, `iat`.

**Multi-tenancy enforcement**: Middleware extracts `school_id` from JWT and injects it into
`context.Context` via a typed key. All repository calls receive this context and scope queries
accordingly.

---

### Database Migrations

**Decision**: `golang-migrate/migrate` with plain SQL files in `migrations/`.

**Naming convention**: `NNN_description.up.sql` / `NNN_description.down.sql`

**Execution**: Migrate container runs before API starts (Docker healthcheck dependency).
In production, migrations are a separate CI step.

---

### Notifications

**Decision**: Push notifications are out of scope for the REST API backend. The API emits
structured notification events; delivery (push/email) is handled by a future async worker.
For MVP, the API will include a notification log table; delivery is deferred.

**Rationale**: Simplicity principle — adding a notification dispatcher now creates coupling
without a clear delivery target (mobile push provider not yet selected).

---

## Dev Tooling (branch `002-makefile-fmt-lint-vulnerability-lefthook`)

### fmt

**Decision**: `go fmt ./...`  
**Rationale**: `go fmt` is the canonical Go formatter; ships with the Go toolchain — no extra install. Idiomatic module-aware form; respects build constraints.  
**Alternatives considered**: `goimports ./...` — handles import ordering too but requires external install. Not justified unless import ordering becomes a recurring code-review pain point.

---

### lint

**Decision**: `golangci-lint run ./...` with a minimal `.golangci.yml`  
**Rationale**: golangci-lint is the de-facto standard aggregator for Go linters. Wraps `go vet`, `staticcheck`, `errcheck`, `gosimple`, and others in a single fast, cached run. A committed `.golangci.yml` pins the enabled linter set so future golangci-lint releases don't introduce noise.  
**Alternatives considered**: `go vet ./...` alone — too limited. `staticcheck ./...` alone — excellent but narrower.

---

### vulnerability

**Decision**: `govulncheck ./...`  
**Rationale**: Official Go vulnerability scanner by the Go Security team (`golang.org/x/vuln`). Traces the call graph and only reports vulnerabilities reachable from actual code paths — very low false-positive rate. Install: `go install golang.org/x/vuln/cmd/govulncheck@latest`.  
**Alternatives considered**: `osv-scanner` — heavier, multi-ecosystem. `nancy` — dependency-list only, no call-graph analysis.

---

### lefthook commit-msg

**Decision**: `lefthook.yml` with a `commit-msg` block using `grep -qE` regex validation  
**Rationale**: Lefthook is a fast, cross-platform git hooks manager written in Go. No Node.js runtime required; version-controlled in the repo; activated with `lefthook install`. The `commit-msg` hook receives the commit message file path as `{1}`.

**Commit message regex**:
```
^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9/_-]+\))?: .+
```

**Alternatives considered**: `commitlint` — requires npm. Manual `.git/hooks/` script — not version-controlled. `husky` — Node.js only.

---

### Installation strategy

Dev tools are **not** added to `go.mod`. A dedicated `make tools` target documents and automates installation. A `make hooks` target runs `lefthook install` to activate the `.git/hooks/commit-msg` symlink.
