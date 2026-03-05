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

var gradleDependencyRe = regexp.MustCompile(`["']([^:'"]+):([^:'"]+):([^'"]+)["']`)

// DependencyReader reads dependencies from build.gradle.
type DependencyReader struct{}

// ReadDependencies parses build.gradle and returns the list of dependencies.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	for _, filename := range []string{"build.gradle", "build.gradle.kts"} {
		path := filepath.Join(repoPath, filename)
		if !fileutil.Exists(path) {
			continue
		}
		content, err := fileutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", filename, err)
		}
		var deps []entities.Dependency
		inDeps := false
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "dependencies {" {
				inDeps = true
				continue
			}
			if inDeps && line == "}" {
				break
			}
			if inDeps {
				if m := gradleDependencyRe.FindStringSubmatch(line); m != nil {
					name := m[1] + ":" + m[2]
					deps = append(deps, entities.NewDependency(name, m[3], "", filename))
				}
			}
		}
		return deps, nil
	}
	return nil, fmt.Errorf("no build.gradle or build.gradle.kts found")
}
