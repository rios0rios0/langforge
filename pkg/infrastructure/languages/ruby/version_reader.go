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
	version, lastReadErr, err := readVersionFromGemspecs(repoPath)
	if err != nil {
		return entities.Version{}, err
	}
	if version != nil {
		return *version, nil
	}

	v, err := readVersionFromFile(repoPath)
	if err != nil {
		return entities.Version{}, err
	}
	if v != nil {
		return *v, nil
	}

	if lastReadErr != nil {
		return entities.Version{}, fmt.Errorf("no version found: %w", lastReadErr)
	}
	return entities.Version{}, errors.New("no version found in .gemspec or VERSION file")
}

func readVersionFromGemspecs(repoPath string) (*entities.Version, error, error) {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading directory %q: %w", repoPath, err)
	}

	var lastReadErr error
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".gemspec") || entry.IsDir() {
			continue
		}
		content, readErr := fileutil.ReadFile(filepath.Join(repoPath, entry.Name()))
		if readErr != nil {
			lastReadErr = fmt.Errorf("reading %s: %w", entry.Name(), readErr)
			continue
		}
		if v := extractGemspecVersion(content); v != nil {
			return v, nil, nil
		}
	}
	return nil, lastReadErr, nil
}

func extractGemspecVersion(content string) *entities.Version {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := gemspecVersionRe.FindStringSubmatch(line); m != nil {
			v, err := entities.NewVersion(m[1])
			if err == nil {
				return &v
			}
		}
	}
	return nil
}

func readVersionFromFile(repoPath string) (*entities.Version, error) {
	versionFile := filepath.Join(repoPath, "VERSION")
	if !fileutil.Exists(versionFile) {
		return nil, nil
	}
	content, err := fileutil.ReadFile(versionFile)
	if err != nil {
		return nil, fmt.Errorf("reading VERSION file: %w", err)
	}
	raw := strings.TrimSpace(content)
	if raw == "" {
		return nil, nil
	}
	v, err := entities.NewVersion(raw)
	if err != nil {
		return nil, err
	}
	return &v, nil
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
