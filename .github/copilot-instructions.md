# Copilot Instructions

## Project Overview

`langforge` is a shared Go library that provides language detection and ecosystem abstractions for Go, Java (Gradle and Maven), Python, Node.js (TypeScript/JavaScript), C#, and Terraform. It exposes interfaces and implementations for:

- Detecting which language/ecosystem a repository uses
- Reading and writing a project's canonical version
- Reading and updating project dependencies

It is used internally by tools such as `autobump` and `autoupdate`.

## Architecture

The project follows **hexagonal architecture** (ports and adapters):

```
pkg/
  domain/
    entities/           # Value objects: Language, Version, Dependency
    repositories/       # Interfaces (ports): LanguageDetector, VersionReader, VersionWriter,
                        #   DependencyReader, DependencyUpdater, BuildValidator, LanguageProvider
  infrastructure/
    languages/          # Implementations (adapters) per language
      golang/           # Go language provider (package name: golang, not go)
      node/
      python/
      java_gradle/
      java_maven/
      csharp/
      terraform/
    registry/           # LanguageRegistry: maps Language → LanguageProvider
  support/
    exec/               # Shell command runner abstraction
    fileutil/           # File read/write helpers
test/
  builders/             # Test builder pattern for domain entities and stubs
  doubles/              # Test doubles (stubs/fakes) for interfaces
```

## Key Interfaces

All language providers satisfy the composite `LanguageProvider` interface defined in `pkg/domain/repositories/language_provider.go`, which combines:

| Interface            | Key Methods                                             |
|---------------------|---------------------------------------------------------|
| `LanguageDetector`  | `Detect(repoPath) (bool, error)`, `Language() Language`, `DetectionFiles() []string` |
| `VersionReader`     | `ReadVersion(repoPath) (Version, error)`, `VersionFiles() []string` |
| `VersionWriter`     | `WriteVersion(repoPath, version) error`, `FilesChanged(repoPath) ([]string, error)` |
| `DependencyReader`  | `ReadDependencies(repoPath) ([]Dependency, error)`      |
| `DependencyUpdater` | `UpdateAll(repoPath) error`, `FilesChanged(repoPath) ([]string, error)`, `Commands() []string` |

Optionally, a provider may also implement `BuildValidator` (`Validate`, `LintCommands`, `BuildCommands`) via `LanguageProviderWithValidation`.

## Adding a New Language Provider

1. Create a new package under `pkg/infrastructure/languages/<language>/`.
2. Add the following files:
   - `detector.go` — implements `LanguageDetector`
   - `version_reader.go` — implements `VersionReader`
   - `version_writer.go` — implements `VersionWriter`
   - `dependency_reader.go` — implements `DependencyReader`
   - `dependency_updater.go` — implements `DependencyUpdater`
   - `provider.go` — `Provider` struct embedding all of the above via pointer composition
3. Register a new `Language` constant in `pkg/domain/entities/language.go`.
4. Wire up the provider in the registry or the consuming application.

### Provider Struct Convention

Use **embedded struct composition** to satisfy the `LanguageProvider` interface:

```go
type Provider struct {
    *Detector
    *VersionReader
    *VersionWriter
    *DependencyReader
    *DependencyUpdater
}

func NewProvider() *Provider {
    return &Provider{
        Detector:          &Detector{},
        VersionReader:     &VersionReader{},
        VersionWriter:     &VersionWriter{},
        DependencyReader:  &DependencyReader{},
        DependencyUpdater: NewDependencyUpdater(exec.NewDefaultRunner()),
    }
}
```

### Package Naming

- Use the ecosystem/tool name as the Go package name (e.g., `java_gradle`, `python`, `node`, `terraform`).
- The Go language implementation uses the package name `golang` (not `go`) to avoid a keyword conflict.

## Domain Entities

### `Language` (`pkg/domain/entities/language.go`)

A `string` type with named constants (e.g., `LanguageGo`, `LanguageNode`, `LanguagePython`).

### `Version` (`pkg/domain/entities/version.go`)

An immutable value object wrapping `github.com/Masterminds/semver/v3`:
- `NewVersion(raw string) (Version, error)` — returns an error for invalid semver
- `MustNewVersion(raw string) Version` — panics on error (use in tests and init code only)

### `Dependency` (`pkg/domain/entities/dependency.go`)

A plain struct with `Name`, `Current`, `Latest`, and `SourceFile` fields:
- `NewDependency(name, current, latest, sourceFile string) Dependency`
- `IsOutdated() bool` — true when `Current != Latest && Latest != ""`

## Testing Patterns

- Tests live alongside the source file they test (e.g., `detector_test.go` next to `detector.go`).
- Use the **builder pattern** from `test/builders/` to construct test entities (`DependencyBuilder`, `VersionBuilder`, `LanguageProviderStubBuilder`).
- Test doubles live in `test/doubles/`.
- Use `github.com/stretchr/testify` for assertions.

## Build and Test Commands

```bash
go build ./...        # build all packages
go test ./...         # run all tests
make lint             # run linters (golangci-lint via pipelines)
make test             # run tests via Makefile
make sast             # run static analysis / security scan
```

The module path is `github.com/rios0rios0/langforge`.

## Code Style

- Follow standard Go conventions (effective Go, `gofmt`).
- Keep interfaces small — each interface in `pkg/domain/repositories/` represents a single responsibility.
- Prefer returning `error` over panicking; only use `Must*` constructors in tests or package-level `init`/`var` blocks.
- All public types and functions must have doc comments.
- Update `CHANGELOG.md` under `[Unreleased]` for every change.
- Commits follow the [rios0rios0 Git Flow conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow).
