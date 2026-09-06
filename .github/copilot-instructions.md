# Copilot Instructions

## Project Overview

`langforge` is a shared Go library that provides language detection and ecosystem abstractions for 11 ecosystems: Go, Node/TypeScript, Python, Dart/Flutter, Java (Gradle and Maven), C#, Ruby, Terraform, Dockerfile, and Pipeline/CI. It exposes interfaces and implementations for:

- Detecting which language/ecosystem a repository uses
- Reading and writing a project's canonical version
- Reading and updating project dependencies

It is used internally by tools such as `autobump` and `autoupdate`.

## Architecture

The project follows **hexagonal architecture** (ports and adapters):

```
pkg/
  domain/
    entities/           # Value objects: Language, Version, Dependency, FileChecker, Classifier
    repositories/       # Interfaces (ports): LanguageDetector, VersionReader, VersionWriter,
                        #   DependencyReader, DependencyUpdater, BuildValidator, RuntimeManager
                        #   Composites: LanguageProvider, LanguageProviderWithValidation, LanguageProviderFull
  infrastructure/
    languages/          # Implementations (adapters) per language
      golang/           # Go language provider (package name: golang, not go)
      node/
      python/
      dart/             # Covers Flutter too (both share pubspec.yaml; IsFlutter picks the toolchain)
      java/             # Not a provider — shared JDK runtime manager for javagradle/javamaven
      toolchain/        # Not a provider — shared BuildValidator/RuntimeManager, configured per language
      javagradle/
      javamaven/
      csharp/
      ruby/
      terraform/
      dockerfile/       # Detection only
      pipeline/         # Detection only
    registry/           # LanguageRegistry: maps Language → LanguageProvider
    versions/           # Version fetchers for latest stable versions from public APIs
  support/
    cmdexec/            # Shell command runner abstraction
    fileutil/           # File read/write helpers, LocalFileChecker
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
| `BuildValidator`    | `Validate(repoPath) error`, `LintCommands() []string`, `BuildCommands() []string` |
| `RuntimeManager`    | `SDKName() string`, `VersionManager() string`, `InstallCommand(version) string`, `CurrentVersion() (string, error)`, `StartCommand() string`, `StopCommand() string` |

Composite interfaces: `LanguageProvider` (first five), `LanguageProviderWithValidation` (adds `BuildValidator`), `LanguageProviderFull` (adds `RuntimeManager`).

## Adding a New Language Provider

1. Create a new package under `pkg/infrastructure/languages/<language>/`.
2. Add the following files:
   - `detector.go` — implements `LanguageDetector`
   - `version_reader.go` — implements `VersionReader`
   - `version_writer.go` — implements `VersionWriter`
   - `dependency_reader.go` — implements `DependencyReader`
   - `dependency_updater.go` — implements `DependencyUpdater`
   - `build_validator.go` — implements `BuildValidator` (optional) by embedding `toolchain.NewBuildValidator` with the language's `toolchain.Commands`
   - `runtime_manager.go` — implements `RuntimeManager` (optional) by embedding `toolchain.NewRuntimeManager` with the language's `toolchain.SDK`
   - `provider.go` — a `NewProvider()` constructor returning a `*repositories.CompositeProvider` wired with those parts
3. Register a new `Language` constant in `pkg/domain/entities/language.go`.
4. Register the provider in `pkg/infrastructure/registry/default_registry.go`. **Order matters:** `Detect` returns the first provider that matches, so a provider with an unambiguous marker file must be registered ahead of one with a weak marker (this is why `dart` precedes `node` — a Flutter web project may keep a `package.json`, but only Dart uses `pubspec.yaml`).

### Provider Composition

Every ecosystem composes the same seven ports the same way, so they all share one type — `repositories.CompositeProvider` — instead of each declaring an identical `Provider` struct. A package supplies only its own parts; unset fields stay nil, and the composite still satisfies the narrower interface (`LanguageProvider` or `LanguageProviderWithValidation`) its callers ask for.

```go
func NewProvider() *repositories.CompositeProvider {
    runner := cmdexec.NewDefaultRunner()
    return &repositories.CompositeProvider{
        LanguageDetector:  &Detector{},
        VersionReader:     &VersionReader{},
        VersionWriter:     &VersionWriter{},
        DependencyReader:  &DependencyReader{},
        DependencyUpdater: NewDependencyUpdater(runner),
        BuildValidator:    NewBuildValidator(runner),  // optional
        RuntimeManager:    NewRuntimeManager(runner),  // optional
    }
}
```

`VersionWriter` and `DependencyUpdater` both declare `FilesChanged`, so the promoted method is ambiguous; `CompositeProvider.FilesChanged` resolves it by merging both results (writer's first, deduplicated), written and tested once instead of per-provider.

### Package Naming

- Use the ecosystem/tool name as the Go package name (e.g., `javagradle`, `python`, `node`, `terraform`).
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

- Build tag: `//go:build unit` on all unit tests.
- External test packages (e.g., `package golang_test`).
- BDD structure with `// given`, `// when`, `// then` comment blocks.
- Parallel execution: `t.Parallel()` at test function and sub-test level.
- Use the **builder pattern** from `test/builders/` to construct test entities (`DependencyBuilder`, `VersionBuilder`, `LanguageProviderStubBuilder`).
- Test doubles live in `test/doubles/`.
- Framework: `github.com/stretchr/testify` (assertions + suites) and `github.com/rios0rios0/testkit`.

## Build and Test Commands

```bash
make setup             # clone/update the pipelines repository (first time)
make lint              # run golangci-lint via pipelines
make test              # run all tests via pipelines
make sast              # run full SAST suite (CodeQL, Semgrep, Trivy, Hadolint, Gitleaks)
go build ./...                                               # build all packages
go test -tags=unit ./...                                     # run all tests directly
go test -tags=unit -v ./pkg/infrastructure/languages/golang  # run tests for a specific package
go test -tags=unit -run TestGoVersionReader ./...             # run a single test by name
```

The module path is `github.com/rios0rios0/langforge`.

## Code Style

- Follow standard Go conventions (effective Go, `gofmt`).
- Keep interfaces small — each interface in `pkg/domain/repositories/` represents a single responsibility.
- Prefer returning `error` over panicking; only use `Must*` constructors in tests or package-level `init`/`var` blocks.
- All public types and functions must have doc comments.
- Add a changelog fragment for every change with `chlog new --kind <Kind> --body "..."`; `CHANGELOG.md` is generated from the fragments and is never edited by hand.
- Commits follow the [rios0rios0 Git Flow conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow).

<!-- chlog:start -->
## Changelog (chlog) — MANDATORY

If the repository you are working in uses chlog (a `.chlog.yaml` or `.chlog.yml`
config file, or a `.changes/` directory, exists at the project root), the
following is binding and ALWAYS applies: whenever you make ANY change, you MUST
create a changelog fragment as part of the same change — automatically, without
being asked, before committing.

- Do NOT edit CHANGELOG.md directly; it is generated from fragments.
- Create the fragment with:
  `chlog new --kind <Kind> --body "<imperative description>"`
- Valid kinds: Added, Changed, Deprecated, Removed, Fixed, Security
- Choose the kind that best matches the change (e.g., new feature → Added,
  bug fix → Fixed, behavior change → Changed, removal → Removed, security fix → Security).
- If the change is backward-INCOMPATIBLE with the public API (a breaking
  change), you MUST add the `--breaking` flag:
  `chlog new --kind <Kind> --breaking --body "<description>"`.
  This is the ONLY thing that triggers a major version bump — the kind alone
  never does (per SemVer, major = incompatible change). When unsure whether a
  change breaks compatibility, ask the user instead of guessing.
- Fragments are YAML files in `.changes/unreleased/`; stage them with your commit.
- `chlog check` fails the build when a fragment is missing — never skip it.
<!-- chlog:end -->
