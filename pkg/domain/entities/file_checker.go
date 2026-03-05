package entities

// FileChecker checks whether a file matching the given path or glob pattern exists.
// For exact paths (e.g. "go.mod"), it checks that specific file.
// For glob patterns (e.g. "*.tf"), it checks if any file matches the pattern.
// This abstraction allows the same detection logic to work with both
// local filesystem access and remote API-based access.
type FileChecker func(pathOrPattern string) (bool, error)
