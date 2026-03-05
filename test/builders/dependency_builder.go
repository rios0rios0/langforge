package builders

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// DependencyBuilder builds Dependency instances using the builder pattern.
type DependencyBuilder struct {
	*testkit.BaseBuilder

	name       string
	current    string
	latest     string
	sourceFile string
}

// NewDependencyBuilder creates a new builder with default values.
func NewDependencyBuilder() *DependencyBuilder {
	return &DependencyBuilder{
		BaseBuilder: testkit.NewBaseBuilder(),
		name:        "test-dependency",
		current:     "1.0.0",
		latest:      "1.1.0",
		sourceFile:  "go.mod",
	}
}

func (b *DependencyBuilder) WithName(name string) *DependencyBuilder {
	b.name = name
	return b
}

func (b *DependencyBuilder) WithCurrent(current string) *DependencyBuilder {
	b.current = current
	return b
}

func (b *DependencyBuilder) WithLatest(latest string) *DependencyBuilder {
	b.latest = latest
	return b
}

func (b *DependencyBuilder) WithSourceFile(sourceFile string) *DependencyBuilder {
	b.sourceFile = sourceFile
	return b
}

func (b *DependencyBuilder) Build() any {
	return entities.NewDependency(b.name, b.current, b.latest, b.sourceFile)
}
