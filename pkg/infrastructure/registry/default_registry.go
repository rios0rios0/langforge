package registry

import (
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/csharp"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dart"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/javagradle"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/javamaven"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/node"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/python"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/ruby"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/terraform"
)

// NewDefaultRegistry creates a LanguageRegistry pre-populated with all built-in language providers.
func NewDefaultRegistry() *LanguageRegistry {
	r := NewLanguageRegistry()
	r.Register(golang.NewProvider())
	// Dart is registered ahead of Node deliberately. Detect returns the first
	// provider that matches, and a Flutter web project often carries a
	// package.json for its JavaScript tooling while nothing but Dart uses a
	// pubspec.yaml — so the more specific marker has to be offered first.
	r.Register(dart.NewProvider())
	r.Register(node.NewProvider())
	r.Register(python.NewProvider())
	r.Register(javagradle.NewProvider())
	r.Register(javamaven.NewProvider())
	r.Register(csharp.NewProvider())
	r.Register(ruby.NewProvider())
	r.Register(terraform.NewProvider())
	return r
}
