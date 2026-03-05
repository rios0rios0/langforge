package registry

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/csharp"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/java_gradle"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/java_maven"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/node"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/python"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/terraform"
)

// NewDefaultRegistry creates a LanguageRegistry pre-populated with all built-in language providers.
func NewDefaultRegistry() *LanguageRegistry {
	r := NewLanguageRegistry()
	r.Register(golang.NewProvider())
	r.Register(node.NewProvider())
	r.Register(python.NewProvider())
	r.Register(java_gradle.NewProvider())
	r.Register(java_maven.NewProvider())
	r.Register(csharp.NewProvider())
	r.Register(terraform.NewProvider())
	return r
}
