package entities

// Language represents a programming language or ecosystem supported by langforge.
type Language string

const (
	LanguageGo         Language = "go"
	LanguageNode       Language = "node"
	LanguagePython     Language = "python"
	LanguageJava       Language = "java"
	LanguageJavaGradle Language = "java_gradle"
	LanguageJavaMaven  Language = "java_maven"
	LanguageRuby       Language = "ruby"
	LanguageCSharp     Language = "csharp"
	LanguageTerraform  Language = "terraform"
	LanguageYAML       Language = "yaml"
	LanguagePipeline   Language = "pipeline"
	LanguageDockerfile Language = "dockerfile"
	LanguageUnknown    Language = "unknown"
)

// String returns the string representation of the language.
func (l Language) String() string {
	return string(l)
}
