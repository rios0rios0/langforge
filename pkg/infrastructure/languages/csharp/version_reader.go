package csharp

import (
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

type csprojPropertyGroup struct {
	Version string `xml:"Version"`
}

type csprojProject struct {
	XMLName        xml.Name              `xml:"Project"`
	PropertyGroups []csprojPropertyGroup `xml:"PropertyGroup"`
}

// VersionReader reads the version from *.csproj.
type VersionReader struct{}

// VersionFiles returns the files inspected for the version.
func (r *VersionReader) VersionFiles() []string {
	return []string{"*.csproj"}
}

// ReadVersion reads the <Version> element from the first *.csproj file.
func (r *VersionReader) ReadVersion(repoPath string) (entities.Version, error) {
	csprojFile, err := findCsprojFile(repoPath)
	if err != nil {
		return entities.Version{}, err
	}
	content, err := fileutil.ReadFile(filepath.Join(repoPath, csprojFile))
	if err != nil {
		return entities.Version{}, fmt.Errorf("reading %s: %w", csprojFile, err)
	}
	var proj csprojProject
	if err := xml.Unmarshal([]byte(content), &proj); err != nil {
		return entities.Version{}, fmt.Errorf("parsing %s: %w", csprojFile, err)
	}
	for _, pg := range proj.PropertyGroups {
		if pg.Version != "" {
			return entities.NewVersion(pg.Version)
		}
	}
	return entities.Version{}, fmt.Errorf("no <Version> element found in %s", csprojFile)
}
