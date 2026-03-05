package csharp

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

type itemGroup struct {
	PackageReferences []packageRef `xml:"PackageReference"`
}

type packageRef struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type csprojDeps struct {
	XMLName    xml.Name    `xml:"Project"`
	ItemGroups []itemGroup `xml:"ItemGroup"`
}

// DependencyReader reads dependencies from *.csproj.
type DependencyReader struct{}

// ReadDependencies parses *.csproj PackageReference elements.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	csprojFile, err := findCsprojFile(repoPath)
	if err != nil {
		return nil, err
	}
	content, err := fileutil.ReadFile(filepath.Join(repoPath, csprojFile))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", csprojFile, err)
	}
	var proj csprojDeps
	if err := xml.Unmarshal([]byte(content), &proj); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", csprojFile, err)
	}
	var deps []entities.Dependency
	for _, ig := range proj.ItemGroups {
		for _, pr := range ig.PackageReferences {
			deps = append(deps, entities.NewDependency(pr.Include, pr.Version, "", csprojFile))
		}
	}
	return deps, nil
}
