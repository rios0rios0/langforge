package javagradle

import (
	"bufio"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var gradleVersionLineRe = regexp.MustCompile(`^(\s*)(version\s*[=:]\s*)["']?[^"'\s]+["']?`)

// VersionWriter updates the version in build.gradle.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	for _, filename := range []string{"build.gradle", "build.gradle.kts"} {
		path := filepath.Join(repoPath, filename)
		if fileutil.Exists(path) {
			return []string{path}, nil
		}
	}
	return nil, errors.New("no build.gradle or build.gradle.kts found")
}

// WriteVersion updates the version field in build.gradle.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	for _, filename := range []string{"build.gradle", "build.gradle.kts"} {
		path := filepath.Join(repoPath, filename)
		if !fileutil.Exists(path) {
			continue
		}
		content, err := fileutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", filename, err)
		}
		var out strings.Builder
		updated := false
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			if !updated && gradleVersionLineRe.MatchString(line) {
				// Preserve the leading whitespace (captured in group 1) and key prefix (group 2)
				line = gradleVersionLineRe.ReplaceAllString(line, fmt.Sprintf(`${1}${2}"%s"`, version.String()))
				updated = true
			}
			out.WriteString(line + "\n")
		}
		if !updated {
			return fmt.Errorf("version field not found in %s", filename)
		}
		return fileutil.WriteFile(path, out.String())
	}
	return errors.New("no build.gradle or build.gradle.kts found")
}
