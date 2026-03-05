package terraform

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects Terraform projects by the presence of *.tf or *.hcl files.
type Detector struct{}

// DetectionFiles returns the files that identify a Terraform project.
func (d *Detector) DetectionFiles() []string {
	return []string{"*.tf", "*.hcl", "versions.tf"}
}

// Detect returns true if any Terraform marker file exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the Terraform language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguageTerraform
}
