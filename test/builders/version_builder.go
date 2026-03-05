package builders

import (
	"github.com/rios0rios0/langforge/pkg/domain/entities"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// VersionBuilder builds Version instances using the builder pattern.
type VersionBuilder struct {
	*testkit.BaseBuilder

	raw string
}

// NewVersionBuilder creates a new builder with a default version.
func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{
		BaseBuilder: testkit.NewBaseBuilder(),
		raw:         "1.0.0",
	}
}

func (b *VersionBuilder) WithRaw(raw string) *VersionBuilder {
	b.raw = raw
	return b
}

func (b *VersionBuilder) Build() any {
	return entities.MustNewVersion(b.raw)
}
