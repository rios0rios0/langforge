package golang

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var goModVersionRe = regexp.MustCompile(`^//\s*version:\s*(.+)$`)

// VersionReader reads the version from a go.mod comment.
// Convention: the first line comment `// version: X.Y.Z` in go.mod holds the version.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"go.mod"}
}

// ReadVersion reads the version from go.mod's version comment.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "go.mod"))
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading go.mod: %w", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if m := goModVersionRe.FindStringSubmatch(line); m != nil {
			return entities.NewVersion(strings.TrimSpace(m[1]))
		}
	}
	return entities.Version{}, fmt.Errorf("no version comment found in go.mod (expected '// version: X.Y.Z')")
}
