# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- changed the Go version to `1.26.6` and updated all module dependencies

### Added

- added a Dart provider covering Flutter as well, since both declare the same `pubspec.yaml`. `IsFlutter` is the single place that distinguishes them — from one manifest, `flutter pub` and `flutter analyze` are selected for a Flutter project and `dart pub`/`dart analyze` for a plain package, because resolving a Flutter project with the bare Dart toolchain cannot satisfy its SDK-sourced packages
- added a line-oriented `pubspec.yaml` version writer rather than a YAML round-trip, because a pubspec routinely records why each dependency was chosen in comments between the entries and re-serialising the document would discard every one of them. Verified against `medhub-tech/frontend-dart`: the rewrite changes the single version line and leaves the file byte-identical elsewhere
- added `BumpBuildNumber`, which carries a Flutter build number across a release and increments it (`1.0.0+1` → `1.1.0+2`), preserving zero padding (`2.10.2+021002` → `2.11.0+021003`) and passing a non-numeric suffix through untouched. Google Play and App Store Connect reject an upload whose build number did not increase, so dropping or freezing it would turn a version bump into a release that cannot ship; a manifest without one keeps none, which is what pub.dev packages want
- added `repositories.CompositeProvider` and `cmdexec.CapturedVersion`, and moved every provider onto them. The nine `Provider` structs were identical, their nine `FilesChanged` methods byte-identical, and seven `RuntimeManager.CurrentVersion` methods differed only in a binary name and a regex — so the new Dart provider could not be written without either duplicating all of it again or extracting it. Each ecosystem now contributes only its own parts, and the merge that resolves the ambiguous `FilesChanged` promotion is written and tested once instead of nine times. `csharp` and `node` keep their own `CurrentVersion`: those genuinely differ, trimming whole output and stripping a `v` prefix rather than matching a pattern
- added a `java` package holding the JDK runtime manager that the Java/Gradle and Java/Maven providers share. The two ecosystems differ in their build tool, not in their SDK, its version manager, or the command that reports its version, so their runtime managers were identical apart from the command that starts an application — which is now the one thing each supplies
- added `.dart` to the extension classifier and registered the provider **ahead of Node** in the default registry. Detection returns the first provider that matches, and a Flutter web project often keeps a `package.json` for its JavaScript tooling while nothing but Dart uses a `pubspec.yaml` — so the more specific marker has to be offered first, and a test pins that ordering

### Removed

- **BREAKING CHANGE:** removed the nine per-ecosystem `Provider` structs, which `repositories.CompositeProvider` replaces. `NewProvider()` still exists in every language package and now returns `*repositories.CompositeProvider`, so registering a provider is unchanged; only code that named a concrete type (`golang.Provider`, `dart.Provider`, …) has to change
- **BREAKING CHANGE:** removed `javagradle.RuntimeManager` and `javamaven.RuntimeManager` along with their constructors, which the shared `java.RuntimeManager` replaces. Both providers report the same SDK as before

## [0.6.11] - 2026-07-14

### Changed

- changed the Go module dependencies to their latest versions

## [0.6.10] - 2026-07-13

### Changed

- changed the Go module dependencies to their latest versions

## [0.6.9] - 2026-07-10

### Changed

- changed the Go version to `1.26.5` and updated all module dependencies

### Security

- replaced `secrets: inherit` with an explicit `CLAUDE_CODE_OAUTH_TOKEN` secret in the Claude workflow callers, following the least-privilege principle

## [0.6.8] - 2026-06-03

### Changed

- changed the Go version to `1.26.4` and updated all module dependencies

## [0.6.7] - 2026-05-19

### Changed

- changed the Go module dependencies to their latest versions
- refreshed `.github/copilot-instructions.md` to add missing `StartCommand()` and `StopCommand()` methods to the `RuntimeManager` interface table

## [0.6.6] - 2026-05-08

### Changed

- changed the Go version to `1.26.3` and updated all module dependencies

## [0.6.5] - 2026-05-01

### Changed

- changed the Go module dependencies to their latest versions

## [0.6.4] - 2026-04-29

### Changed

- changed the Go module dependencies to their latest versions

## [0.6.3] - 2026-04-28

### Changed

- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to document Ruby provider, `BuildValidator`/`RuntimeManager` interfaces, `LanguageProviderFull` composite, `versions/` package, and corrected package names (`javagradle`/`javamaven`, `cmdexec`)

## [0.6.2] - 2026-04-16

### Changed

- changed the Go module dependencies to their latest versions

## [0.6.1] - 2026-04-15

### Changed

- changed the Go version to `1.26.2` and updated all module dependencies

## [0.6.0] - 2026-04-14

### Added

- added Ruby language provider with detection (`Gemfile`, `*.gemspec`), version reading/writing (`.gemspec` and `VERSION` file), dependency reading, `bundle update` integration, `BuildValidator`, and `RuntimeManager`

## [0.5.0] - 2026-03-23

### Added

- added `BuildValidator` implementations for all languages (Go, Node, Python, Java/Gradle, Java/Maven, C#, Terraform)
- added `LanguageProviderFull` composite interface combining validation and runtime management
- added `RunOutput` method to `cmdexec.Runner` interface for capturing command stdout
- added `RuntimeManager` implementations for all languages with version detection and manager commands
- added `RuntimeManager` interface for SDK/runtime detection and installation commands

## [0.4.0] - 2026-03-17

### Added

- added `httptest`-based unit tests for `fetchJSON`, `fetchEndOfLifeLatest`, and all version fetcher selection/filtering logic

### Changed

- changed `fetchJSON` and `fetchEndOfLifeLatest` to accept an injectable `*http.Client` parameter for testability

### Fixed

- fixed `FetchLatestJavaVersion` returning 404 by switching from the non-existent `/api/java.json` endpoint to `/api/amazon-corretto.json` on endoflife.date

## [0.3.1] - 2026-03-14

### Changed

- changed the Go module dependencies to their latest versions

## [0.3.0] - 2026-03-12

### Added

- added `NewFileChecker`, `IsGlobPattern`, and `ExtractExtension` utilities in `pkg/support/fileutil/` to enable building remote-compatible `FileChecker` instances from custom callbacks, eliminating the need for consumers to reimplement glob-vs-exact dispatch
- added `VersionFetcher` type and `FetchLatestGoVersion`, `FetchLatestNodeVersion`, `FetchLatestPythonVersion`, `FetchLatestJavaVersion`, `FetchLatestTerraformVersion` helpers in `pkg/infrastructure/versions/` for fetching latest stable versions from public APIs

### Changed

- changed `fetchJSON` to limit response body size via `io.LimitReader` to prevent memory blowups from misbehaving endpoints
- changed `fetchJSON` to respect caller-provided context deadlines instead of hard-coding an HTTP client timeout
- changed `isActiveEOL` to compare dates in UTC and treat the EOL date as inclusive
- changed the Go version to `1.26.1` and updated all module dependencies

## [0.2.0] - 2026-03-09

### Added

- added `DockerfileDetector` for detecting repositories containing Dockerfile files
- added `LanguagePipeline` and `LanguageDockerfile` constants for CI/CD and container ecosystem detection
- added `PipelineDetector` for detecting repositories with GitHub Actions workflows and Azure DevOps pipeline templates

## [0.1.0] - 2026-03-06

### Added

- added `*.hcl` to Terraform detector's detection files
- added `DetectWith` standalone function that runs detection against any `FileChecker`
- added `DetectWithChecker` and `DetectAllWithChecker` methods to `LanguageRegistry` for remote-compatible detection
- added `FileChecker` function type to enable pluggable file existence checks (local filesystem or remote API)
- added `LanguageJava` and `LanguageYAML` constants to `Language` type
- added `LanguageRegistry` with auto-detection logic
- added `LocalFileChecker` constructor that creates a `FileChecker` backed by the local filesystem
- added `NewDefaultRegistry` convenience constructor pre-populated with all built-in language providers
- added `requirements.txt` to Python detector's detection files
- added file extension classifier with `ClassifyFileByExtension` and `ClassifyFilesByExtension` functions
- added implementations for Go, Node/TypeScript, Python, Java/Gradle, Java/Maven, C#, and Terraform
- added initial library scaffold with language abstraction layer
- added language detection, version reading/writing, and dependency management interfaces
- added support packages for file utilities and command execution

### Changed

- changed the Go version to `1.26.0` and updated all module dependencies
- refactored all 7 language detectors to use `DetectWith` internally, eliminating duplicated detection logic

### Fixed

- fixed `FilesChanged` ambiguity on all `Provider` structs by adding explicit disambiguation between `VersionWriter` and `DependencyUpdater`
- fixed `funcorder` findings by reordering constructors before methods in all 7 language provider files
- fixed `godoclint` findings by using proper doc comment link syntax
- fixed `govet` shadow findings by eliminating variable shadowing in 8 dependency reader/writer files
- fixed `mnd` finding by extracting magic number constant in Go dependency reader
- fixed `noctx` finding by using context-aware command execution in `cmdexec` package
- fixed `staticcheck` and `revive` package naming findings by renaming `java_gradle` to `javagradle`, `java_maven` to `javamaven`, and `exec` to `cmdexec`
- fixed `unparam` finding by removing always-nil error return from `resolveRefTagLine` in Terraform dependency updater
