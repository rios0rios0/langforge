package entities

// Language represents a programming language or ecosystem supported by langforge.
type Language string

const (
	LanguageGo         Language = "go"
	LanguageNode       Language = "node"
	LanguagePython     Language = "python"
	LanguageJavaGradle Language = "java_gradle"
	LanguageJavaMaven  Language = "java_maven"
	LanguageCSharp     Language = "csharp"
	LanguageTerraform  Language = "terraform"
	LanguageUnknown    Language = "unknown"
)

// String returns the string representation of the language.
func (l Language) String() string {
	return string(l)
}
