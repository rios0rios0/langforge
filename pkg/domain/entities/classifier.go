package entities

import "path/filepath"

// extensionToLanguage maps file extensions to their canonical Language.
//
//nolint:gochecknoglobals // read-only lookup table used as a constant
var extensionToLanguage = map[string]Language{
	".go":   LanguageGo,
	".js":   LanguageNode,
	".ts":   LanguageNode,
	".jsx":  LanguageNode,
	".tsx":  LanguageNode,
	".mjs":  LanguageNode,
	".cjs":  LanguageNode,
	".py":   LanguagePython,
	".java": LanguageJava,
	".cs":   LanguageCSharp,
	".tf":   LanguageTerraform,
	".hcl":  LanguageTerraform,
	".yaml": LanguageYAML,
	".yml":  LanguageYAML,
}

// ClassifyFileByExtension returns the Language for a file path based on its extension.
// Returns LanguageUnknown if the extension is not recognized.
func ClassifyFileByExtension(filePath string) Language {
	ext := filepath.Ext(filePath)
	if lang, ok := extensionToLanguage[ext]; ok {
		return lang
	}
	return LanguageUnknown
}

// ClassifyFilesByExtension returns the unique set of languages for the given file paths.
// Files with unrecognized extensions are excluded from the result.
func ClassifyFilesByExtension(paths []string) []Language {
	seen := make(map[Language]bool)
	var languages []Language

	for _, p := range paths {
		lang := ClassifyFileByExtension(p)
		if lang != LanguageUnknown && !seen[lang] {
			seen[lang] = true
			languages = append(languages, lang)
		}
	}

	return languages
}
