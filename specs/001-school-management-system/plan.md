# Implementation Plan: Makefile fmt/lint/vulnerability + Lefthook commit-msg

**Branch**: `002-makefile-fmt-lint-vulnerability-lefthook` | **Date**: 2026-04-24 | **Spec**: `specs/001-school-management-system/spec.md`
**Input**: Add `make fmt`, `make lint`, `make vulnerability` targets and Lefthook `commit-msg` hook with conventional commit validation.

## Summary

Add three developer workflow targets to the existing Makefile (`fmt`, `lint`, `vulnerability`) using standard Go tooling, and introduce a `lefthook.yml` that enforces the conventional commit message format `type(scope)?: description` via the `commit-msg` git hook.

## Technical Context

**Language/Version**: Go 1.26.2  
**Primary Dependencies**: golangci-lint (external), govulncheck (`golang.org/x/vuln`), lefthook  
**Storage**: N/A  
**Testing**: N/A  
**Target Platform**: macOS developer machine / Linux CI  
**Project Type**: Dev tooling configuration (Makefile + git hooks)  
**Performance Goals**: N/A  
**Constraints**: No new runtime dependencies; tooling must be installable via `go install` or Homebrew  
**Scale/Scope**: Developer workflow (single repo, small team)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean Architecture | ✅ PASS | No application code modified; tooling/config only |
| II. Clean Code | ✅ PASS | Makefile targets are single-purpose, named clearly |
| III. Simplicity (YAGNI) | ✅ PASS | `go fmt ./...` uses stdlib; only well-justified external tools added |
| IV. Go Idiomatic Patterns | ✅ PASS | `go fmt ./...` and `govulncheck` are canonical Go toolchain commands |
| V. Observability | ✅ PASS | N/A for tooling |

No gate violations. No complexity tracking needed.

**Post-design re-check**: Same verdict. Config files and Makefile additions introduce zero application-layer complexity.

## Project Structure

### Documentation (this feature)

```text
specs/001-school-management-system/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output (config schemas)
├── quickstart.md        ← Phase 1 output
└── contracts/
    └── commit-message.md  ← Phase 1 output
```

### Source Code (repository root)

```text
Makefile          # Updated: add fmt, lint, vulnerability, tools, hooks targets
lefthook.yml      # New: commit-msg hook configuration
.golangci.yml     # New: golangci-lint linter selection
```

**Structure Decision**: Single-project layout. All changes are root-level config/tooling files and additions to the existing Makefile. No new directories under `internal/` or `cmd/` are required.

## Complexity Tracking

> No constitution violations — table omitted per template instructions.
