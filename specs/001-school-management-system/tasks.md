# Tasks: Makefile fmt/lint/vulnerability + Lefthook commit-msg

**Input**: Design documents from `/specs/001-school-management-system/` (branch `002-makefile-fmt-lint-vulnerability-lefthook`)  
**Plan**: `plan.md` (branch `002-*`)  
**Tests**: None — constitution prohibits unit tests; tooling validated by manual invocation.  
**Organization**: Single-phase feature. All tasks are independent prerequisites to application development.

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- No user stories — this is infrastructure/tooling work

---

## Phase 1: Setup

**Purpose**: Verify working state before making changes.

- [x] T001 Confirm `Makefile` exists at repo root and lists current targets (`run`, `build`, `migrate`, `docker-up`, `docker-down`, `tidy`)

---

## Phase 2: Makefile Targets

**Purpose**: Add developer workflow commands to `Makefile`.

- [x] T002 Update `Makefile` — extend `.PHONY` to include `fmt lint vulnerability tools hooks` and add the five new targets:
  - `fmt`: `go fmt ./...`
  - `lint`: `golangci-lint run ./...`
  - `vulnerability`: `govulncheck ./...`
  - `tools`: `go install golang.org/x/vuln/cmd/govulncheck@latest` + install instructions for golangci-lint and lefthook
  - `hooks`: `lefthook install`

**Checkpoint**: `make fmt` formats Go source; `make lint` runs linters; `make vulnerability` scans deps.

---

## Phase 3: Config Files

**Purpose**: Create config files required by the new Makefile targets.

- [x] T003 [P] Create `.golangci.yml` at repo root — enable linters: errcheck, gosimple, govet, staticcheck, unused; set `check-blank: true` for errcheck; exclude `_test.go` files from errcheck; limit `max-same-issues: 3`
- [x] T004 [P] Create `lefthook.yml` at repo root — `commit-msg` hook with a `validate-message` command that reads `cat "{1}"` and validates against regex `^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9/_-]+\))?: .+`; on failure, print the expected format and examples before exiting 1

**Checkpoint**: `.golangci.yml` and `lefthook.yml` exist and are valid YAML.

---

## Phase 4: Activation & Validation

**Purpose**: Activate git hooks and verify end-to-end.

- [x] T005 Install dev tools: run `make tools` to install `govulncheck`; install `golangci-lint` via Homebrew (`brew install golangci-lint`) or the official install script; install lefthook (`brew install lefthook` or `go install github.com/evilmartians/lefthook@latest`)
- [x] T006 Run `make hooks` (`lefthook install`) to activate the `commit-msg` hook in `.git/hooks/`
- [x] T007 [P] Verify `make fmt` runs without error on the current codebase
- [x] T008 [P] Verify `make lint` runs without error (or with only acceptable violations)
- [x] T009 [P] Verify `make vulnerability` runs and reports no known vulnerabilities
- [x] T010 Verify the commit-msg hook rejects `bad message` and accepts `feat(auth): add login`

**Checkpoint**: All four make targets run clean. Lefthook rejects non-conventional commit messages.

---

## Dependencies & Execution Order

- **Phase 1 (Setup)**: No dependencies
- **Phase 2 (Makefile)**: Depends on Phase 1 verification
- **Phase 3 (Config files)**: Can run in parallel with Phase 2 (different files)
- **Phase 4 (Activation)**: Depends on Phases 2 + 3 complete

### Parallel Opportunities

```bash
# Phase 3 — run these together (different files):
Task T003: Create .golangci.yml
Task T004: Create lefthook.yml

# Phase 4 verification — run these together after tools installed:
Task T007: make fmt
Task T008: make lint
Task T009: make vulnerability
```

---

## Implementation Strategy

### Single-increment delivery

1. T001: Verify starting state
2. T002: Update Makefile (the core change)
3. T003 + T004 in parallel: Create config files
4. T005 + T006: Activate tooling
5. T007–T010: Validate everything works

Total: 10 tasks, ~30 minutes of work.

---

## Notes

- [P] tasks operate on different files and can run concurrently
- No user story labels — this is dev tooling, not feature work
- After T006, every `git commit` will validate the message format automatically
- `make fmt` should be run before every commit; add to pre-commit hook in a future iteration if desired
- `govulncheck` requires Go module with a valid `go.sum`; ensure `go mod tidy` was run first (already in Makefile as `make tidy`)
