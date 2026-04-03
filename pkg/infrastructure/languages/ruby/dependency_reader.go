package ruby

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var gemDeclRe = regexp.MustCompile(`^\s*gem\s+["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`)

// DependencyReader reads dependencies from Gemfile or .gemspec.
type DependencyReader struct{}

// ReadDependencies reads dependencies from Gemfile or .gemspec.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	gemfilePath := filepath.Join(repoPath, "Gemfile")
	if fileutil.Exists(gemfilePath) {
		return readGemfile(gemfilePath)
	}

	// Try .gemspec file
	entries, err := os.ReadDir(repoPath)
	if err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".gemspec") && !entry.IsDir() {
				return readGemspec(filepath.Join(repoPath, entry.Name()))
			}
		}
	}

	return nil, errors.New("no Gemfile or .gemspec found")
}

func readGemfile(path string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading Gemfile: %w", err)
	}

	var deps []entities.Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if m := gemDeclRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			version := ""
			if len(m) > 2 {
				version = m[2]
			}
			deps = append(deps, entities.NewDependency(name, version, "", "Gemfile"))
		}
	}

	return deps, nil
}

var addDepRe = regexp.MustCompile(
	`\.\s*add(?:_runtime)?_dependency\s*\(?\s*["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`,
)

func readGemspec(path string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading .gemspec: %w", err)
	}

	var deps []entities.Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if m := addDepRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			version := ""
			if len(m) > 2 {
				version = m[2]
			}
			deps = append(deps, entities.NewDependency(name, version, "", filepath.Base(path)))
		}
	}

	return deps, nil
}
