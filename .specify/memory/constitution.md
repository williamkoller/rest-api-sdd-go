<!--
SYNC IMPACT REPORT
==================
Version change: (template) → 1.0.0
Modified principles: N/A (initial ratification from template)
Added sections:
  - I. Clean Architecture
  - II. Clean Code
  - III. Simplicity (YAGNI)
  - IV. Go Idiomatic Patterns
  - V. Observability
  - Project Structure
  - Development Workflow
  - Governance
Removed sections: N/A
Templates requiring updates:
  - .specify/templates/plan-template.md ✅ — Constitution Check gates align (no test gates)
  - .specify/templates/spec-template.md ✅ — No test-related mandatory sections added
  - .specify/templates/tasks-template.md ✅ — Tests marked OPTIONAL; no unit tests required
Follow-up TODOs: None. All placeholders resolved.
-->

# Go REST API Constitution

## Core Principles

### I. Clean Architecture

The codebase MUST follow strict layer separation with a single direction of dependency:
`domain` ← `application` ← `infrastructure` ← `transport (HTTP)`.

- The `domain` layer contains entities and business rules. It MUST NOT import any other
  internal layer or external framework.
- The `application` layer contains use cases (interactors). It MUST depend only on `domain`
  interfaces, never on infrastructure or transport details.
- The `infrastructure` layer contains database adapters, external service clients, and
  repository implementations. It MUST implement interfaces defined in `domain` or `application`.
- The `transport` layer (HTTP handlers, routers) MUST depend only on `application` use cases
  via interfaces. It MUST NOT contain business logic.
- Dependency inversion MUST be achieved through Go interfaces defined at the consumer side
  (application or domain), not at the implementation side.

**Rationale**: Decoupled layers allow the persistence mechanism, HTTP framework, and external
integrations to change without touching business logic.

### II. Clean Code

All code MUST be readable, expressive, and self-documenting.

- Functions and methods MUST do one thing only. If a function needs a comment to explain
  what it does, it MUST be refactored into smaller named functions.
- Names MUST be intentional: variables, functions, and types express their purpose without
  abbreviation. Single-letter names are only acceptable for trivial loop indices.
- Functions MUST be short. A function exceeding ~30 lines SHOULD be decomposed.
- No dead code, commented-out code, or TODO comments may be merged to the main branch
  without an associated issue.
- Magic numbers and magic strings MUST be extracted to named constants.
- Public API surface MUST be minimal — only export what callers need.

**Rationale**: Code is read far more often than it is written. Clean, expressive code reduces
cognitive load and prevents bugs caused by misunderstanding.

### III. Simplicity (YAGNI)

The project MUST remain simple and avoid speculative complexity.

- No abstraction, pattern, or indirection MUST be introduced until it solves a concrete,
  present problem. Hypothetical future requirements do NOT justify added complexity.
- Dependencies (third-party modules) MUST be justified. Prefer the Go standard library.
  Each new dependency requires an explicit rationale.
- No unit tests. The project is simple enough that correctness is validated through
  manual testing, integration checks, and clear domain logic. Unit test infrastructure
  MUST NOT be added without explicit team decision and constitution amendment.
- Database migrations MUST be plain SQL files. ORM layers are prohibited unless justified.
- Configuration MUST come from environment variables. No complex config frameworks.

**Rationale**: Complexity is the primary source of bugs and maintenance cost. Simple, direct
code is easier to understand, change, and debug than an over-engineered codebase.

### IV. Go Idiomatic Patterns

Code MUST follow Effective Go conventions and Go Code Review Comments.

- Errors MUST be returned, never ignored (use `_` only for provably irrelevant returns).
- Errors MUST be wrapped with context using `fmt.Errorf("operation: %w", err)` at each
  layer boundary.
- Sentinel errors and custom error types MUST be used for errors that callers need to
  distinguish programmatically.
- Interfaces MUST be small (1–3 methods). Large interfaces indicate a design smell.
- `context.Context` MUST be the first parameter of any function that performs I/O or can
  be cancelled, and MUST be propagated — never stored in structs.
- Goroutines MUST have explicit lifecycle management: every goroutine started MUST have a
  clear, documented exit condition.
- Struct embedding is acceptable for code reuse, but MUST NOT be used to fake inheritance.

**Rationale**: Idiomatic Go is predictable. Following established conventions reduces
surprises, aids code review, and leverages Go tooling effectively.

### V. Observability

The service MUST be observable in production without requiring a debugger.

- Structured logging MUST use `log/slog` with consistent field names across all handlers.
- Every HTTP request MUST log: method, path, status code, and duration.
- Every error returned to the client MUST be logged with the full internal error chain
  (use `%+v` or `errors.Unwrap` chain) at `ERROR` level.
- No `fmt.Println` or `log.Printf` (unstructured logging) in production code paths.
- Health check endpoint (`GET /health`) MUST be implemented and return service status.

**Rationale**: Invisible services are unmanageable. Consistent structured logs are the
minimum observability floor for any HTTP service.

## Project Structure

Standard Clean Architecture layout for this Go REST API:

```text
cmd/
└── api/
    └── main.go              # Entry point — wires dependencies, starts server

internal/
├── domain/                  # Entities, value objects, repository interfaces
│   ├── entity.go
│   └── repository.go
├── application/             # Use cases / interactors
│   └── usecase.go
├── infrastructure/          # DB adapters, external clients
│   └── postgres/
│       └── repository.go
└── transport/               # HTTP handlers, routers, middleware
    └── http/
        ├── handler.go
        └── router.go

migrations/                  # Plain SQL migration files
config/                      # Environment-based configuration loading
```

The `internal/` boundary is enforced by Go's package visibility rules. No package outside
this module MAY import `internal/`.

## Development Workflow

- Every feature starts with a specification (`/speckit-specify`) and a plan (`/speckit-plan`).
- Tasks are generated from the plan (`/speckit-tasks`) and implemented sequentially or in
  parallel as marked.
- Code review MUST verify compliance with all five principles before merge.
- All work happens on feature branches; main branch MUST remain deployable at all times.
- Database changes MUST include a migration file. Schema changes without a migration are
  prohibited.

## Governance

This constitution supersedes all other documented practices for this project. All agents,
contributors, and code reviewers MUST treat these principles as non-negotiable constraints.

Amendment procedure:
1. Propose amendment with rationale and impact on existing code.
2. Update this document following semantic versioning rules (see Version line).
3. Run `/speckit-constitution` to propagate changes to all dependent templates and docs.
4. Commit with message: `docs: amend constitution to vX.Y.Z (<summary>)`.

Compliance review: Every pull request MUST include a "Constitution Check" section in the
plan or review confirming all five principles are satisfied. Complexity violations MUST be
explicitly justified in the plan's Complexity Tracking table.

**Version**: 1.0.0 | **Ratified**: 2026-04-23 | **Last Amended**: 2026-04-23
