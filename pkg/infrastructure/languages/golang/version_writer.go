package golang

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var goModVersionLineRe = regexp.MustCompile(`(?m)^(// version:).*$`)

// VersionWriter updates the version comment in go.mod.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	return []string{filepath.Join(repoPath, "go.mod")}, nil
}

// WriteVersion updates the `// version: X.Y.Z` comment in go.mod.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, "go.mod")
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading go.mod: %w", err)
	}
	if !goModVersionLineRe.MatchString(content) {
		// Prepend the version comment before the module line
		content = fmt.Sprintf("// version: %s\n%s", version.String(), content)
	} else {
		content = goModVersionLineRe.ReplaceAllStringFunc(content, func(_ string) string {
			return fmt.Sprintf("// version: %s", version.String())
		})
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return fileutil.WriteFile(path, content)
}
