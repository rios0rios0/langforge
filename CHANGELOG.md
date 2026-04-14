# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
