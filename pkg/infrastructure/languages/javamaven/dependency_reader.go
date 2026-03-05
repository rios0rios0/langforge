package javamaven

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

type pomDep struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type pomDeps struct {
	XMLName      xml.Name `xml:"project"`
	Dependencies []pomDep `xml:"dependencies>dependency"`
}

// DependencyReader reads dependencies from pom.xml.
type DependencyReader struct{}

// ReadDependencies parses pom.xml and returns the list of dependencies.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, "pom.xml"))
	if err != nil {
		return nil, fmt.Errorf("reading pom.xml: %w", err)
	}
	var proj pomDeps
	if err = xml.Unmarshal([]byte(content), &proj); err != nil {
		return nil, fmt.Errorf("parsing pom.xml: %w", err)
	}
	var deps []entities.Dependency
	for _, d := range proj.Dependencies {
		name := d.GroupID + ":" + d.ArtifactID
		deps = append(deps, entities.NewDependency(name, d.Version, "", "pom.xml"))
	}
	return deps, nil
}
