package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Terraform projects by the presence of *.tf files.
type Detector struct{}

// DetectionFiles returns the files that identify a Terraform project.
func (d *Detector) DetectionFiles() []string {
	return []string{"*.tf", "versions.tf"}
}

// Detect returns true if any *.tf file exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	if fileutil.Exists(filepath.Join(repoPath, "versions.tf")) {
		return true, nil
	}
	matches, err := fileutil.GlobFiles(repoPath, "*.tf")
	if err != nil {
		return false, fmt.Errorf("globbing *.tf: %w", err)
	}
	return len(matches) > 0, nil
}

// Language returns the Terraform language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageTerraform
}
