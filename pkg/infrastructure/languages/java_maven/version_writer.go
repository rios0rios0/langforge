package java_maven

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var pomVersionRe = regexp.MustCompile(`(?m)(<version>)[^<]*(</version>)`)

// VersionWriter updates the version in pom.xml.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	return []string{filepath.Join(repoPath, "pom.xml")}, nil
}

// WriteVersion updates the first <version> element in pom.xml.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, "pom.xml")
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading pom.xml: %w", err)
	}
	updated := false
	newContent := pomVersionRe.ReplaceAllStringFunc(content, func(match string) string {
		if updated {
			return match
		}
		updated = true
		return fmt.Sprintf("<version>%s</version>", version.String())
	})
	if !updated {
		return fmt.Errorf("no <version> element found in pom.xml")
	}
	return fileutil.WriteFile(path, newContent)
}
