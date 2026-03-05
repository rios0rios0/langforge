package terraform

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

var (
	terraformVersionRe    = regexp.MustCompile(`required_version\s*=\s*["']([^"']+)["']`)
	terraformVersionNumRe = regexp.MustCompile(`(\d+\.\d+(?:\.\d+)*)`)
)

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
			// Extract the bare version number from the constraint (e.g. ">= 1.5.0" → "1.5.0")
			if numMatch := terraformVersionNumRe.FindString(m[1]); numMatch != "" {
				return entities.NewVersion(numMatch)
			}
		}
	}
	return entities.Version{}, errors.New("no required_version found in versions.tf")
}
