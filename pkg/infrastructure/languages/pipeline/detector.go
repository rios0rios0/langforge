package pipeline

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Detector detects CI/CD pipeline configurations by the presence of
// GitHub Actions workflows or Azure DevOps pipeline templates.
type Detector struct{}

// DetectionFiles returns the files that identify a CI/CD pipeline project.
func (d *Detector) DetectionFiles() []string {
	return []string{
		".github/workflows/*.yaml",
		".github/workflows/*.yml",
		".azure-pipelines.yml",
		"azure-pipelines.yml",
		"azure-devops/*.yaml",
		"azure-devops/*/*.yaml",
		"azure-devops/*/*/*.yaml",
	}
}

// Detect returns true if any pipeline configuration file exists in repoPath.
func (d *Detector) Detect(repoPath string) (bool, error) {
	return repositories.DetectWith(d, fileutil.LocalFileChecker(repoPath))
}

// Language returns the pipeline language identifier.
func (d *Detector) Language() entities.Language {
	return entities.LanguagePipeline
}
