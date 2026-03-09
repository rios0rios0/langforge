# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`langforge` is a shared Go library (module `github.com/rios0rios0/langforge`) that provides language detection, version management, and dependency abstractions for 9 ecosystems: Go, Node/TypeScript, Python, Java Gradle, Java Maven, C#, Terraform, Dockerfile, and Pipeline/CI. It is consumed by tools like `autobump` and `autoupdate`. This is a library — there is no CLI or `main` package.

## Commands

```bash
make setup             # clone/update the pipelines repository (first time)
make lint              # run golangci-lint via pipelines
make test              # run all tests via pipelines
make sast              # run full SAST suite (CodeQL, Semgrep, Trivy, Hadolint, Gitleaks)
go build ./...         # build all packages
go test ./...          # run all tests directly
go test -v ./pkg/infrastructure/languages/golang  # run tests for a specific package
go test -run TestGoVersionReader ./...             # run a single test by name
```

Always use `make lint` / `make test` / `make sast` instead of calling tool binaries directly.

## Architecture

Hexagonal Architecture (Ports & Adapters) with DDD:

- **`pkg/domain/entities/`** — Value objects: `Language` (enum), `Version` (semver wrapper), `Dependency`, `FileChecker` (pluggable file existence), `Classifier` (extension-based)
- **`pkg/domain/repositories/`** — Interface contracts (ports): `LanguageDetector`, `VersionReader`, `VersionWriter`, `DependencyReader`, `DependencyUpdater`, `BuildValidator`. The composite `LanguageProvider` combines the first five.
- **`pkg/infrastructure/languages/<ecosystem>/`** — One package per ecosystem implementing the domain interfaces. Each has: `detector.go`, `version_reader.go`, `version_writer.go`, `dependency_reader.go`, `dependency_updater.go`, `provider.go`. Dockerfile and Pipeline packages implement only `Detector`.
- **`pkg/infrastructure/registry/`** — `LanguageRegistry` maps `Language` → `LanguageProvider`; `NewDefaultRegistry()` wires all built-in providers.
- **`pkg/support/`** — Cross-cutting utilities: `cmdexec/` (shell command runner), `fileutil/` (file I/O helpers, `LocalFileChecker`).
- **`test/builders/`** — Builder pattern for test data (`VersionBuilder`, `DependencyBuilder`, `LanguageProviderStubBuilder`).
- **`test/doubles/`** — Test doubles (stubs) implementing domain interfaces.

### Provider pattern

Each language `Provider` uses embedded struct composition to satisfy `LanguageProvider`:

```go
type Provider struct {
    *Detector
    *VersionReader
    *VersionWriter
    *DependencyReader
    *DependencyUpdater
}
```

The Go language implementation uses package name `golang` (not `go`) to avoid the keyword conflict.

## Adding a new language provider

1. Create package under `pkg/infrastructure/languages/<name>/`
2. Implement `Detector`, `VersionReader`, `VersionWriter`, `DependencyReader`, `DependencyUpdater` (or a subset for detection-only providers)
3. Create `Provider` struct with embedded composition and a `NewProvider()` constructor
4. Add a `Language` constant in `pkg/domain/entities/language.go`
5. Register in `pkg/infrastructure/registry/default_registry.go`

## Testing conventions

- Build tag: `//go:build unit` on all unit tests
- External test packages (e.g., `package golang_test`)
- BDD structure with `// given`, `// when`, `// then` comment blocks
- Parallel execution: `t.Parallel()` at test function and sub-test level
- Use builders from `test/builders/` and stubs from `test/doubles/`
- Framework: `github.com/stretchr/testify` (assertions + suites) and `github.com/rios0rios0/testkit`

## Key dependencies

- `github.com/Masterminds/semver/v3` — semantic version parsing
- `github.com/stretchr/testify` — test assertions and suites
- `github.com/rios0rios0/testkit` — shared test utilities
