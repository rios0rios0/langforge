package python

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var pyprojectVersionLineRe = regexp.MustCompile(`^(version\s*=\s*)["']?[^"'\s]+["']?`)

// VersionWriter updates the version in pyproject.toml.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	return []string{filepath.Join(repoPath, "pyproject.toml")}, nil
}

// WriteVersion updates the version field in pyproject.toml.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, "pyproject.toml")
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading pyproject.toml: %w", err)
	}

	var out strings.Builder
	inProject := false
	updated := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "[project]" {
			inProject = true
		} else if inProject && strings.HasPrefix(trimmed, "[") {
			inProject = false
		}
		if inProject && !updated && pyprojectVersionLineRe.MatchString(trimmed) {
			line = pyprojectVersionLineRe.ReplaceAllString(trimmed, fmt.Sprintf(`${1}"%s"`, version.String()))
			updated = true
		}
		out.WriteString(line + "\n")
	}
	if !updated {
		return fmt.Errorf("version field not found in [project] section of pyproject.toml")
	}
	return fileutil.WriteFile(path, out.String())
}
