# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-03-06

### Added

- added `*.hcl` to Terraform detector's detection files
- added `DetectWithChecker` and `DetectAllWithChecker` methods to `LanguageRegistry` for remote-compatible detection
- added `DetectWith` standalone function that runs detection against any `FileChecker`
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
