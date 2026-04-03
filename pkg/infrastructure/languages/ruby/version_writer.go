package ruby

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var gemspecVersionLineRe = regexp.MustCompile(`(\.version\s*=\s*)["'][^"']+["']`)

// VersionWriter updates the version in a .gemspec file.
type VersionWriter struct{}

// FilesChanged returns the list of files that will be modified.
func (w *VersionWriter) FilesChanged(repoPath string) ([]string, error) {
	gemspecPath, err := findGemspecFile(repoPath)
	if err == nil {
		return []string{gemspecPath}, nil
	}
	versionFile := filepath.Join(repoPath, "VERSION")
	if fileutil.Exists(versionFile) {
		return []string{versionFile}, nil
	}
	return nil, fmt.Errorf("no *.gemspec or VERSION file found in %q", repoPath)
}

// WriteVersion updates the version field in the .gemspec or VERSION file.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	gemspecPath, err := findGemspecFile(repoPath)
	if err == nil {
		return writeGemspecVersion(gemspecPath, version)
	}

	// Fall back to VERSION file
	versionFile := filepath.Join(repoPath, "VERSION")
	if fileutil.Exists(versionFile) {
		return fileutil.WriteFile(versionFile, version.String()+"\n")
	}

	return fmt.Errorf("no *.gemspec or VERSION file found in %q", repoPath)
}

func writeGemspecVersion(gemspecPath string, version entities.Version) error {
	content, err := fileutil.ReadFile(gemspecPath)
	if err != nil {
		return fmt.Errorf("reading .gemspec: %w", err)
	}

	var out strings.Builder
	updated := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if !updated && gemspecVersionLineRe.MatchString(line) {
			line = gemspecVersionLineRe.ReplaceAllString(
				line, fmt.Sprintf(`${1}"%s"`, version.String()),
			)
			updated = true
		}
		out.WriteString(line + "\n")
	}
	if !updated {
		return fmt.Errorf("version specification not found in %s", gemspecPath)
	}
	return fileutil.WriteFile(gemspecPath, out.String())
}
