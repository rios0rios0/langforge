package ruby

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var gemspecVersionRe = regexp.MustCompile(`\.version\s*=\s*["']([^"']+)["']`)

// VersionReader reads the version from a .gemspec file or VERSION file.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"*.gemspec", "VERSION"}
}

// ReadVersion reads the version from .gemspec or VERSION file.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	// Try .gemspec files first
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading directory %q: %w", repoPath, err)
	}

	var lastReadErr error
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".gemspec") && !entry.IsDir() {
			content, readErr := fileutil.ReadFile(filepath.Join(repoPath, entry.Name()))
			if readErr != nil {
				lastReadErr = fmt.Errorf("reading %s: %w", entry.Name(), readErr)
				continue
			}
			scanner := bufio.NewScanner(strings.NewReader(content))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if m := gemspecVersionRe.FindStringSubmatch(line); m != nil {
					return entities.NewVersion(m[1])
				}
			}
		}
	}

	// Try VERSION file
	versionFile := filepath.Join(repoPath, "VERSION")
	if fileutil.Exists(versionFile) {
		content, readErr := fileutil.ReadFile(versionFile)
		if readErr != nil {
			return entities.Version{}, fmt.Errorf("reading VERSION file: %w", readErr)
		}
		version := strings.TrimSpace(content)
		if version != "" {
			return entities.NewVersion(version)
		}
	}

	if lastReadErr != nil {
		return entities.Version{}, fmt.Errorf("no version found: %w", lastReadErr)
	}
	return entities.Version{}, errors.New("no version found in .gemspec or VERSION file")
}

// findGemspecFile returns the first .gemspec file found in the given directory.
func findGemspecFile(repoPath string) (string, error) {
	matches, err := fileutil.GlobFiles(repoPath, "*.gemspec")
	if err != nil {
		return "", fmt.Errorf("globbing *.gemspec: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no *.gemspec file found in %q", repoPath)
	}
	return matches[0], nil
}
