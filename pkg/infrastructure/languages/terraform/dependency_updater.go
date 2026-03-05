package terraform

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/support/exec"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

var tfRefTagRe = regexp.MustCompile(`\?ref=([^\s"]+)`)

// DependencyUpdater performs custom ref-tag resolution for Terraform modules.
type DependencyUpdater struct {
	runner exec.Runner
}

// NewDependencyUpdater creates a DependencyUpdater with the default runner.
func NewDependencyUpdater(runner exec.Runner) *DependencyUpdater {
	return &DependencyUpdater{runner: runner}
}

// Commands returns the shell commands that UpdateAll runs.
func (u *DependencyUpdater) Commands() []string {
	return []string{
		"terraform init -upgrade",
	}
}

// FilesChanged returns the files modified by an update.
func (u *DependencyUpdater) FilesChanged(repoPath string) ([]string, error) {
	matches, err := fileutil.GlobFiles(repoPath, "*.tf")
	if err != nil {
		return nil, err
	}
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m
	}
	return result, nil
}

// UpdateAll runs terraform init -upgrade.
func (u *DependencyUpdater) UpdateAll(repoPath string) error {
	return u.runner.Run(repoPath, "terraform", "init", "-upgrade")
}

// UpdateRefTags updates ?ref=<tag> references in Terraform module sources.
func UpdateRefTags(repoPath string, resolver func(source string) (string, error)) error {
	matches, err := fileutil.GlobFiles(repoPath, "*.tf")
	if err != nil {
		return fmt.Errorf("globbing *.tf: %w", err)
	}
	for _, path := range matches {
		if err := updateRefTagsInFile(path, resolver); err != nil {
			return err
		}
	}
	return nil
}

func updateRefTagsInFile(
	path string,
	resolver func(source string) (string, error),
) error {
	content, err := fileutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	var out strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if tfRefTagRe.MatchString(line) {
			newLine, err := resolveRefTagLine(line, resolver)
			if err != nil {
				return err
			}
			line = newLine
		}
		out.WriteString(line + "\n")
	}
	return fileutil.WriteFile(path, out.String())
}

func resolveRefTagLine(
	line string,
	resolver func(source string) (string, error),
) (string, error) {
	return tfRefTagRe.ReplaceAllStringFunc(line, func(match string) string {
		newTag, err := resolver(match)
		if err != nil {
			return match
		}
		return fmt.Sprintf("?ref=%s", newTag)
	}), nil
}

// ProviderVersion represents a resolved Terraform provider version.
type ProviderVersion struct {
	Source  string
	Version string
}

// ParseProviderVersions extracts provider sources and versions from a versions.tf file.
func ParseProviderVersions(repoPath string) ([]ProviderVersion, error) {
	deps, err := (&DependencyReader{}).ReadDependencies(repoPath)
	if err != nil {
		return nil, err
	}
	result := make([]ProviderVersion, 0, len(deps))
	for _, d := range deps {
		result = append(result, ProviderVersion{
			Source:  d.Name,
			Version: d.Current,
		})
	}
	return result, nil
}
