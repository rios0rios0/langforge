package dart

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// VersionReader reads the version from pubspec.yaml.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{Pubspec}
}

// ReadVersion reads the version field from pubspec.yaml.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, Pubspec))
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading %s: %w", Pubspec, err)
	}
	raw, ok := ExtractVersion(content)
	if !ok {
		// An application always carries one, but a package in a pub workspace
		// legitimately omits it, so this is a plain condition rather than a
		// malformed-file error.
		return entities.Version{}, errors.New("no version field in " + Pubspec)
	}
	return entities.NewVersion(raw)
}

// ExtractVersion returns the raw version value of a pubspec.yaml, and whether
// the manifest declared one. The build suffix is preserved: semver treats it as
// metadata, and Flutter treats it as the store build number.
func ExtractVersion(content string) (string, bool) {
	match := VersionLineRe.FindStringSubmatch(content)
	if match == nil {
		return "", false
	}
	return match[2], true
}
