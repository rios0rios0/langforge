package csharp

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var csprojVersionRe = regexp.MustCompile(`(?m)(<Version>)[^<]*(</Version>)`)

// VersionWriter updates the version in *.csproj.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	csprojFile, err := findCsprojFile(repoPath)
	if err != nil {
		return nil, err
	}
	return []string{filepath.Join(repoPath, csprojFile)}, nil
}

// WriteVersion updates the <Version> element in the *.csproj file.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	csprojFile, err := findCsprojFile(repoPath)
	if err != nil {
		return err
	}
	path := filepath.Join(repoPath, csprojFile)
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", csprojFile, err)
	}
	updated := false
	newContent := csprojVersionRe.ReplaceAllStringFunc(content, func(match string) string {
		if updated {
			return match
		}
		updated = true
		return fmt.Sprintf("<Version>%s</Version>", version.String())
	})
	if !updated {
		return fmt.Errorf("no <Version> element found in %s", csprojFile)
	}
	return fileutil.WriteFile(path, newContent)
}
