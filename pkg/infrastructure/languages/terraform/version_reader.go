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

var terraformVersionRe = regexp.MustCompile(`required_version\s*=\s*["']([^"']+)["']`)

// VersionReader reads the required_version from versions.tf.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"versions.tf"}
}

// ReadVersion reads the required_version constraint from versions.tf.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	path := filepath.Join(repoPath, "versions.tf")
	if !fileutil.Exists(path) {
		return entities.Version{}, fmt.Errorf("versions.tf not found in %q", repoPath)
	}
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading versions.tf: %w", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := terraformVersionRe.FindStringSubmatch(line); m != nil {
			ver := strings.TrimLeft(m[1], ">= ~!")
			return entities.NewVersion(strings.TrimSpace(ver))
		}
	}
	return entities.Version{}, fmt.Errorf("no required_version found in versions.tf")
}
