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

var gradleVersionRe = regexp.MustCompile(`(?m)^version\s*[=:]\s*["']?([^"'\s]+)["']?`)

// VersionReader reads the version from build.gradle.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"build.gradle", "build.gradle.kts"}
}

// ReadVersion reads the version from build.gradle or build.gradle.kts.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	for _, filename := range r.VersionFiles() {
		path := filepath.Join(repoPath, filename)
		if !fileutil.Exists(path) {
			continue
		}
		content, err := fileutil.ReadFile(path)
		if err != nil {
			return entities.Version{}, fmt.Errorf("reading %s: %w", filename, err)
		}
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if m := gradleVersionRe.FindStringSubmatch(line); m != nil {
				return entities.NewVersion(m[1])
			}
		}
	}
	return entities.Version{}, fmt.Errorf("no version found in build.gradle or build.gradle.kts")
}
