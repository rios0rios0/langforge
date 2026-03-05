package java_maven

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

type pomProject struct {
	XMLName xml.Name `xml:"project"`
	Version string   `xml:"version"`
}

// VersionReader reads the version from pom.xml.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"pom.xml"}
}

// ReadVersion reads the version from pom.xml.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "pom.xml"))
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading pom.xml: %w", err)
	}
	var proj pomProject
	if err := xml.Unmarshal([]byte(content), &proj); err != nil {
		return entities.Version{}, fmt.Errorf("parsing pom.xml: %w", err)
	}
	if proj.Version == "" {
		return entities.Version{}, fmt.Errorf("no version element found in pom.xml")
	}
	return entities.NewVersion(proj.Version)
}
