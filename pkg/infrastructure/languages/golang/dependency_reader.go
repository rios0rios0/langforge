package golang

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

const minRequireFields = 2

// DependencyReader reads dependencies from go.mod.
type DependencyReader struct{}

// ReadDependencies parses go.mod require blocks and returns the list of dependencies.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("reading go.mod: %w", err)
	}

	var deps []entities.Dependency
	inRequire := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if after, ok := strings.CutPrefix(line, "require "); ok {
			// single-line require
			line = after
			inRequire = false
		} else if !inRequire {
			continue
		}
		// Remove inline comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		parts := strings.Fields(line)
		if len(parts) < minRequireFields {
			continue
		}
		deps = append(deps, entities.NewDependency(parts[0], parts[1], "", "go.mod"))
	}
	return deps, nil
}
