package dart

import (
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// VersionWriter updates the version in pubspec.yaml.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	path := filepath.Join(repoPath, Pubspec)
	if !fileutil.Exists(path) {
		return nil, fmt.Errorf("no %s found in %q", Pubspec, repoPath)
	}
	return []string{path}, nil
}

// WriteVersion updates the version field in pubspec.yaml, preserving the build
// number and every other byte of the file.
//
// The rewrite is a targeted replacement of one line rather than a YAML
// marshal round-trip, because a pubspec routinely documents why each dependency
// was chosen in comments between the entries, and re-serialising the document
// would discard all of them.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	path := filepath.Join(repoPath, Pubspec)
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", Pubspec, err)
	}

	updated, ok := ReplaceVersion(content, version.String())
	if !ok {
		return fmt.Errorf("version specification not found in %s", path)
	}
	return fileutil.WriteFile(path, updated)
}

// ReplaceVersion rewrites the version line of a pubspec.yaml to newSemver,
// carrying the existing build number forward, and reports whether a version
// line was found. Only the first match is replaced.
func ReplaceVersion(content, newSemver string) (string, bool) {
	loc := pubspecVersionLineRe.FindStringSubmatchIndex(content)
	if loc == nil {
		return content, false
	}
	current := content[loc[4]:loc[5]]
	replacement := content[loc[2]:loc[3]] + BumpBuildNumber(current, newSemver) + content[loc[6]:loc[7]]
	return content[:loc[0]] + replacement + content[loc[1]:], true
}
