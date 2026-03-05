package entities

// Dependency represents a project dependency with its name, current version,
// latest available version, and the source file it was found in.
type Dependency struct {
	Name       string
	Current    string
	Latest     string
	SourceFile string
}

// NewDependency creates a new Dependency.
func NewDependency(name, current, latest, sourceFile string) Dependency {
	return Dependency{
		Name:       name,
		Current:    current,
		Latest:     latest,
		SourceFile: sourceFile,
	}
}

// IsOutdated returns true if the current version differs from the latest.
func (d Dependency) IsOutdated() bool {
	return d.Current != d.Latest && d.Latest != ""
}
