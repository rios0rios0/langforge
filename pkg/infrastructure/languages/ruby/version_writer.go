package ruby

import (
	"bufio"
	"errors"
	"fmt"
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
	if err != nil {
		return nil, err
	}
	return []string{gemspecPath}, nil
}

// WriteVersion updates the version field in the .gemspec file.
func (w *VersionWriter) WriteVersion(repoPath string, version entities.Version) error {
	gemspecPath, err := findGemspecFile(repoPath)
	if err != nil {
		return err
	}

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
		return errors.New("version specification not found in .gemspec file")
	}
	return fileutil.WriteFile(gemspecPath, out.String())
}
