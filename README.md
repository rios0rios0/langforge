<h1 align="center">langforge</h1>
<p align="center">
    <a href="https://github.com/rios0rios0/langforge/releases/latest">
        <img src="https://img.shields.io/github/release/rios0rios0/langforge.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/rios0rios0/langforge/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/rios0rios0/langforge.svg?style=for-the-badge&logo=github" alt="License"/></a>
    <a href="https://github.com/rios0rios0/langforge/actions/workflows/default.yaml">
        <img src="https://img.shields.io/github/actions/workflow/status/rios0rios0/langforge/default.yaml?branch=main&style=for-the-badge&logo=github" alt="Build Status"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_langforge">
        <img src="https://img.shields.io/sonar/coverage/rios0rios0_langforge?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Coverage"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_langforge">
        <img src="https://img.shields.io/sonar/quality_gate/rios0rios0_langforge?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Quality Gate"/></a>
    <a href="https://www.bestpractices.dev/projects/12091">
        <img src="https://img.shields.io/cii/level/12091?style=for-the-badge&logo=opensourceinitiative" alt="OpenSSF Best Practices"/></a>
</p>

A shared Go library providing language detection and ecosystem abstractions for multiple programming languages. Includes version reading/writing, dependency management, and build validation. Used by [autobump](https://github.com/rios0rios0/autobump) and [autoupdate](https://github.com/rios0rios0/autoupdate).

## Features

- **Language Detection**: Automatically detect project languages via manifest files and file extensions
- **Version Management**: Read and write semantic versions across all supported ecosystems
- **Dependency Management**: Read dependency manifests and run native update commands (go get, npm update, pip, etc.)
- **Build Validation**: Run ecosystem-specific lint and build commands to verify project health
- **Remote-Compatible Detection**: Pluggable `FileChecker` abstraction enables detection via GitHub/GitLab APIs without local filesystem access
- **Registry Pattern**: `LanguageRegistry` with auto-detection and provider lookup for polyglot projects
- **File Classification**: Extension-based file classifier for fast, deterministic language identification

## Supported Ecosystems

| Ecosystem       | Detection Files                                  | Version File     |
|-----------------|--------------------------------------------------|------------------|
| Go              | `go.mod`                                         | `go.mod`         |
| Node/TypeScript | `package.json`, `yarn.lock`, `pnpm-lock.yaml`    | `package.json`   |
| Python          | `pyproject.toml`, `setup.py`, `requirements.txt` | `pyproject.toml` |
| Java (Gradle)   | `build.gradle`, `build.gradle.kts`               | `build.gradle`   |
| Java (Maven)    | `pom.xml`                                        | `pom.xml`        |
| C#              | `.csproj`, `.sln`                                | `.csproj`        |
| Terraform       | `.tf`, `.hcl`                                    | `*.tf`           |

## Installation

```bash
go get github.com/rios0rios0/langforge
```

## Usage

```go
import (
    "github.com/rios0rios0/langforge/pkg/infrastructure/registry"
)

// Create registry with all built-in providers
reg := registry.NewDefaultRegistry()

// Auto-detect language from a local project directory
provider, err := reg.Detect("/path/to/repo")

// Or detect using a custom file checker (e.g., GitHub API)
provider, err := reg.DetectWithChecker(myRemoteFileChecker)

// Detect all languages in a polyglot project
providers, err := reg.DetectAllWithChecker(myRemoteFileChecker)

// Read version, dependencies, update, validate
version, err := provider.ReadVersion("/path/to/repo")
deps, err := provider.ReadDependencies("/path/to/repo")
```

## Architecture

```
langforge/
├── pkg/
│   ├── domain/
│   │   ├── entities/         # Language, Version, Dependency, FileChecker, Classifier
│   │   └── repositories/     # LanguageProvider, LanguageDetector, VersionReader/Writer,
│   │                         # DependencyReader, DependencyUpdater, BuildValidator
│   ├── infrastructure/
│   │   ├── languages/        # Go, Node, Python, Java (Gradle/Maven), C#, Terraform
│   │   └── registry/         # LanguageRegistry with auto-detection and default setup
│   └── support/
│       ├── cmdexec/          # Shell command execution wrapper
│       └── fileutil/         # File I/O helpers, LocalFileChecker
└── test/
    ├── builders/             # Test data builders (Version, Dependency, Provider)
    └── doubles/              # Test doubles (stubs)
```

## Provider Interfaces

The library uses Go interface composition through the `LanguageProvider` contract:

- **`LanguageDetector`**: Detects if a language is present in a directory
- **`VersionReader`** / **`VersionWriter`**: Read and write semantic versions in ecosystem-specific files
- **`DependencyReader`**: Parse dependency manifests and return structured dependency lists
- **`DependencyUpdater`**: Run native ecosystem update commands
- **`BuildValidator`**: Validate project builds and linting

Each language provider composes these interfaces into a single `Provider` struct, enabling independent testing and reuse.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

See [LICENSE](LICENSE) file for details.
