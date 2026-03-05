package terraform

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
	tfProviderSourceRe  = regexp.MustCompile(`source\s*=\s*["']([^"']+)["']`)
	tfProviderVersionRe = regexp.MustCompile(`version\s*=\s*["']([^"']+)["']`)
)

// DependencyReader reads provider dependencies from versions.tf or *.tf files.
type DependencyReader struct{}

// ReadDependencies reads required_providers from versions.tf.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	path := filepath.Join(repoPath, "versions.tf")
	if !fileutil.Exists(path) {
		return nil, fmt.Errorf("versions.tf not found")
	}
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading versions.tf: %w", err)
	}

	var deps []entities.Dependency
	inRequiredProviders := false
	inProviderBlock := false
	currentProvider := ""
	currentSource := ""
	currentVersion := ""

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "required_providers") && strings.HasSuffix(line, "{") {
			inRequiredProviders = true
			continue
		}
		if inRequiredProviders && strings.HasSuffix(line, "{") {
			inProviderBlock = true
			currentProvider = strings.TrimSuffix(strings.TrimSpace(line), " {")
			currentSource = ""
			currentVersion = ""
			continue
		}
		if inProviderBlock && line == "}" {
			if currentProvider != "" {
				name := currentSource
				if name == "" {
					name = currentProvider
				}
				deps = append(deps, entities.NewDependency(name, currentVersion, "", "versions.tf"))
			}
			inProviderBlock = false
			continue
		}
		if inRequiredProviders && line == "}" {
			break
		}
		if inProviderBlock {
			if m := tfProviderSourceRe.FindStringSubmatch(line); m != nil {
				currentSource = m[1]
			}
			if m := tfProviderVersionRe.FindStringSubmatch(line); m != nil {
				currentVersion = m[1]
			}
		}
	}
	return deps, nil
}
