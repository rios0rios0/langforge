# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- added `LanguageJava` and `LanguageYAML` constants to `Language` type
- added file extension classifier with `ClassifyFileByExtension` and `ClassifyFilesByExtension` functions
- added `NewDefaultRegistry` convenience constructor pre-populated with all built-in language providers
- added `FileChecker` function type to enable pluggable file existence checks (local filesystem or remote API)
- added `LocalFileChecker` constructor that creates a `FileChecker` backed by the local filesystem
- added `DetectWith` standalone function that runs detection against any `FileChecker`
- added `DetectWithChecker` and `DetectAllWithChecker` methods to `LanguageRegistry` for remote-compatible detection
- added `requirements.txt` to Python detector's detection files
- added `*.hcl` to Terraform detector's detection files

### Fixed

- fixed `FilesChanged` ambiguity on all `Provider` structs by adding explicit disambiguation between `VersionWriter` and `DependencyUpdater`

### Changed

- changed the Go version to `1.26.0` and updated all module dependencies
- refactored all 7 language detectors to use `DetectWith` internally, eliminating duplicated detection logic

### Added
- Initial library scaffold with language abstraction layer
- Language detection, version reading/writing, and dependency management interfaces
- Implementations for Go, Node/TypeScript, Python, Java/Gradle, Java/Maven, C#, and Terraform
- LanguageRegistry with auto-detection logic
- Support packages for file utilities and command execution
