package java_gradle

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var gradleVersionLineRe = regexp.MustCompile(`^(version\s*[=:]\s*)["']?[^"'\s]+["']?`)

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
	return nil, fmt.Errorf("no build.gradle or build.gradle.kts found")
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
			trimmed := strings.TrimSpace(line)
			if !updated && gradleVersionLineRe.MatchString(trimmed) {
				line = gradleVersionLineRe.ReplaceAllString(trimmed, fmt.Sprintf(`${1}"%s"`, version.String()))
				updated = true
			}
			out.WriteString(line + "\n")
		}
		if !updated {
			return fmt.Errorf("version field not found in %s", filename)
		}
		return fileutil.WriteFile(path, out.String())
	}
	return fmt.Errorf("no build.gradle or build.gradle.kts found")
}
