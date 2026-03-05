package terraform

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var terraformVersionLineRe = regexp.MustCompile(`(required_version\s*=\s*")([^"]+)(")`)

// VersionWriter updates the required_version in versions.tf.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	return []string{filepath.Join(repoPath, "versions.tf")}, nil
}

// WriteVersion updates the required_version in versions.tf.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, "versions.tf")
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading versions.tf: %w", err)
	}
	updated := false
	newContent := terraformVersionLineRe.ReplaceAllStringFunc(content, func(match string) string {
		if updated {
			return match
		}
		updated = true
		return fmt.Sprintf(`required_version = ">= %s"`, version.String())
	})
	if !updated {
		return errors.New("no required_version found in versions.tf")
	}
	return fileutil.WriteFile(path, newContent)
}
