# Commit Message Contract

**Feature**: `002-makefile-fmt-lint-vulnerability-lefthook`  
**Enforcement**: Lefthook `commit-msg` hook (`lefthook.yml`)

---

## Format

```
type(scope)?: description
```

## Types

| Type | When to use |
|------|-------------|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, whitespace (no logic change) |
| `refactor` | Code restructuring (no feature or bug change) |
| `perf` | Performance improvement |
| `test` | Adding or correcting tests |
| `build` | Build system, Makefile, Docker changes |
| `ci` | CI/CD pipeline changes |
| `chore` | Maintenance, dependency bumps, tooling |
| `revert` | Reverts a previous commit |

## Scope

Optional. Identifies the module or area changed. Use lowercase with `/`, `-`, or `_` as separators.

**Examples**: `auth`, `financial`, `financial/invoice`, `agenda`, `db`, `deps`, `docker`

## Description

- Imperative mood: "add" not "added" or "adds"
- Lowercase first letter
- No period at the end
- Max ~72 characters

## Validation Regex

```
^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9/_-]+\))?: .+
```

## Valid Examples

```
feat(auth): add JWT refresh token endpoint
fix(financial): handle nil pointer in invoice generation
docs: update API setup instructions
build(docker): add redis healthcheck to compose
chore(deps): bump golang.org/x/vuln to v1.1.0
refactor(agenda): extract item validation to domain layer
ci: add golangci-lint step to GitHub Actions
```

## Invalid Examples

```
Added JWT refresh             ← missing type prefix
feat: Added JWT refresh       ← "Added" is past tense, uppercase
feat(AUTH): add something     ← scope must be lowercase
feat(auth) add something      ← missing colon after scope
```
