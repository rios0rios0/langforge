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

var (
	pyprojectVersionRe = regexp.MustCompile(`(?m)^version\s*=\s*["']?([^"'\s]+)["']?`)
	initVersionRe      = regexp.MustCompile(`__version__\s*=\s*["']([^"']+)["']`)
)

// VersionReader reads the version from pyproject.toml or __init__.py.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"pyproject.toml", "__init__.py"}
}

// ReadVersion reads the version from pyproject.toml first, then __init__.py.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	pyprojectPath := filepath.Join(repoPath, "pyproject.toml")
	if fileutil.Exists(pyprojectPath) {
		content, err := fileutil.ReadFile(pyprojectPath)
		if err != nil {
			return entities.Version{}, fmt.Errorf("reading pyproject.toml: %w", err)
		}
		if v := extractVersionFromPyproject(content); v != "" {
			return entities.NewVersion(v)
		}
	}
	initPath := filepath.Join(repoPath, "__init__.py")
	if fileutil.Exists(initPath) {
		content, err := fileutil.ReadFile(initPath)
		if err != nil {
			return entities.Version{}, fmt.Errorf("reading __init__.py: %w", err)
		}
		if m := initVersionRe.FindStringSubmatch(content); m != nil {
			return entities.NewVersion(m[1])
		}
	}
	return entities.Version{}, fmt.Errorf("no version found in pyproject.toml or __init__.py")
}

func extractVersionFromPyproject(content string) string {
	inProject := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "[project]" {
			inProject = true
			continue
		}
		if inProject && strings.HasPrefix(line, "[") {
			inProject = false
			continue
		}
		if inProject {
			if m := pyprojectVersionRe.FindStringSubmatch(line); m != nil {
				return m[1]
			}
		}
	}
	return ""
}
