package python

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// DependencyReader reads dependencies from requirements.txt or pyproject.toml.
type DependencyReader struct{}

// ReadDependencies reads dependencies from requirements.txt or pyproject.toml.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	reqPath := filepath.Join(repoPath, "requirements.txt")
	if fileutil.Exists(reqPath) {
		return readRequirementsTxt(reqPath)
	}
	pyprojectPath := filepath.Join(repoPath, "pyproject.toml")
	if fileutil.Exists(pyprojectPath) {
		return readPyprojectDeps(pyprojectPath)
	}
	return nil, fmt.Errorf("no requirements.txt or pyproject.toml found")
}

func readRequirementsTxt(path string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading requirements.txt: %w", err)
	}
	var deps []entities.Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		for _, sep := range []string{"==", ">=", "<=", "!=", "~=", ">"} {
			if idx := strings.Index(line, sep); idx != -1 {
				name := strings.TrimSpace(line[:idx])
				ver := strings.TrimSpace(line[idx+len(sep):])
				deps = append(deps, entities.NewDependency(name, sep+ver, "", "requirements.txt"))
				goto nextLine
			}
		}
		deps = append(deps, entities.NewDependency(line, "", "", "requirements.txt"))
	nextLine:
	}
	return deps, nil
}

func readPyprojectDeps(path string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading pyproject.toml: %w", err)
	}
	var deps []entities.Dependency
	inDeps := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "dependencies = [" || line == "dependencies=[" {
			inDeps = true
			continue
		}
		if inDeps && line == "]" {
			break
		}
		if inDeps {
			line = strings.Trim(line, `"',`)
			if line == "" {
				continue
			}
			deps = append(deps, entities.NewDependency(line, "", "", "pyproject.toml"))
		}
	}
	return deps, nil
}
