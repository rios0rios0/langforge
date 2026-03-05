package csharp

import (
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects C# projects by the presence of *.sln or *.csproj files.
type Detector struct{}

// DetectionFiles returns the files that identify a C# project.
func (d *Detector) DetectionFiles() []string {
	return []string{"*.sln", "*.csproj"}
}

// Detect returns true if any *.sln or *.csproj file exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the C# language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageCSharp
}

// findCsprojFile returns the first *.csproj file found in repoPath.
func findCsprojFile(repoPath string) (string, error) {
	matches, err := fileutil.GlobFiles(repoPath, "*.csproj")
	if err != nil {
		return "", fmt.Errorf("globbing *.csproj: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no *.csproj file found in %q", repoPath)
	}
	return filepath.Base(matches[0]), nil
}
