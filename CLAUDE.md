# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`langforge` is a shared Go library (module `github.com/rios0rios0/langforge`) that provides language detection, version management, and dependency abstractions for 11 ecosystems: Go, Node/TypeScript, Python, Dart/Flutter, Java Gradle, Java Maven, C#, Ruby, Terraform, Dockerfile, and Pipeline/CI. It is consumed by tools like `autobump` and `autoupdate`. This is a library — there is no CLI or `main` package.

## Commands

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

Prefer `make lint` / `make test` / `make sast` for CI-equivalent validation. Direct `go test -tags=unit` is fine for quick local iteration.

## Architecture

Hexagonal Architecture (Ports & Adapters) with DDD:

- **`pkg/domain/entities/`** — Value objects: `Language` (enum), `Version` (semver wrapper), `Dependency`, `FileChecker` (pluggable file existence), `Classifier` (extension-based)
- **`pkg/domain/repositories/`** — Interface contracts (ports): `LanguageDetector`, `VersionReader`, `VersionWriter`, `DependencyReader`, `DependencyUpdater`, `BuildValidator`, `RuntimeManager`. Composite interfaces: `LanguageProvider` (first five), `LanguageProviderWithValidation` (adds `BuildValidator`), `LanguageProviderFull` (adds `RuntimeManager`).
- **`pkg/infrastructure/languages/<ecosystem>/`** — One package per ecosystem implementing the domain interfaces. Full providers have: `detector.go`, `version_reader.go`, `version_writer.go`, `dependency_reader.go`, `dependency_updater.go`, `build_validator.go`, `runtime_manager.go`, `provider.go`. Dockerfile and Pipeline packages implement only `Detector`.
- **`pkg/infrastructure/languages/toolchain/`** — Not a provider: the `BuildValidator` and `RuntimeManager` implementations every ecosystem composes, configured per language with `toolchain.Commands` (lint and build command specs) and `toolchain.SDK` (name, version manager, install/start/stop commands, version command and capture pattern). A language's `build_validator.go` and `runtime_manager.go` embed these and supply only their own data.
- **`pkg/infrastructure/registry/`** — `LanguageRegistry` maps `Language` → `LanguageProvider`; `NewDefaultRegistry()` wires all built-in providers.
- **`pkg/infrastructure/versions/`** — Version fetchers for retrieving latest stable versions from public APIs (endoflife.date, etc.).
- **`pkg/support/`** — Cross-cutting utilities: `cmdexec/` (shell command runner), `fileutil/` (file I/O helpers, `LocalFileChecker`).
- **`test/builders/`** — Builder pattern for test data (`VersionBuilder`, `DependencyBuilder`, `LanguageProviderStubBuilder`).
- **`test/doubles/`** — Test doubles (stubs) implementing domain interfaces.

### Provider pattern

Every ecosystem composes the same seven ports the same way, so they share one composite type — `repositories.CompositeProvider` — instead of each declaring an identical struct. A language package supplies only its own parts:

```go
func NewProvider() *repositories.CompositeProvider {
    runner := cmdexec.NewDefaultRunner()
    return &repositories.CompositeProvider{
        LanguageDetector:  &Detector{},
        VersionReader:     &VersionReader{},
        VersionWriter:     &VersionWriter{},
        DependencyReader:  &DependencyReader{},
        DependencyUpdater: NewDependencyUpdater(runner),
        BuildValidator:    NewBuildValidator(runner),  // optional — implements LanguageProviderWithValidation
        RuntimeManager:    NewRuntimeManager(runner),  // optional — implements LanguageProviderFull
    }
}
```

Unset fields stay nil, and the composite still satisfies the narrower interface (`LanguageProvider` or `LanguageProviderWithValidation`) its callers ask for. `VersionWriter` and `DependencyUpdater` both declare `FilesChanged`, so the promoted method is ambiguous; `CompositeProvider.FilesChanged` resolves it by merging both results, written and tested once.

The Go language implementation uses package name `golang` (not `go`) to avoid the keyword conflict.

The `java` package is not a provider — it holds the JDK runtime manager that the `javagradle` and `javamaven` providers share, since the two differ in their build tool rather than in their SDK.

The `dart` package covers Flutter too — both declare the same `pubspec.yaml`, so they are one provider rather than two, and `dart.IsFlutter` is the single place that decides which toolchain (`dart` or `flutter`) a given repository is driven with.

## Adding a new language provider

1. Create package under `pkg/infrastructure/languages/<name>/`
2. Implement `Detector`, `VersionReader`, `VersionWriter`, `DependencyReader`, `DependencyUpdater`, and optionally `BuildValidator` and `RuntimeManager` (or a subset for detection-only providers). For the last two, embed `toolchain.NewBuildValidator` / `toolchain.NewRuntimeManager` and supply the language's `toolchain.Commands` and `toolchain.SDK` instead of re-implementing the ports
3. Create a `NewProvider()` constructor returning a `*repositories.CompositeProvider` wired with those parts
4. Add a `Language` constant in `pkg/domain/entities/language.go`
5. Register in `pkg/infrastructure/registry/default_registry.go` — **order matters**: `Detect` returns the first provider that matches, so a provider whose marker file is unambiguous must be registered ahead of one whose marker is weak (this is why `dart` precedes `node`)

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
