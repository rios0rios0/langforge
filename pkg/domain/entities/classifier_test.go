//go:build unit

package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
)

func TestClassifyFileByExtension(t *testing.T) {
	t.Parallel()

	t.Run("should return LanguageGo when file has .go extension", func(t *testing.T) {
		// given
		filePath := "main.go"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageGo, result)
	})

	t.Run("should return LanguageNode when file has .js extension", func(t *testing.T) {
		// given
		filePath := "index.js"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageNode, result)
	})

	t.Run("should return LanguageNode when file has .tsx extension", func(t *testing.T) {
		// given
		filePath := "component.tsx"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageNode, result)
	})

	t.Run("should return LanguagePython when file has .py extension", func(t *testing.T) {
		// given
		filePath := "script.py"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguagePython, result)
	})

	t.Run("should return LanguageJava when file has .java extension", func(t *testing.T) {
		// given
		filePath := "Main.java"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageJava, result)
	})

	t.Run("should return LanguageCSharp when file has .cs extension", func(t *testing.T) {
		// given
		filePath := "Program.cs"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageCSharp, result)
	})

	t.Run("should return LanguageTerraform when file has .tf extension", func(t *testing.T) {
		// given
		filePath := "main.tf"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageTerraform, result)
	})

	t.Run("should return LanguageYAML when file has .yaml extension", func(t *testing.T) {
		// given
		filePath := "config.yaml"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageYAML, result)
	})

	t.Run("should return LanguageYAML when file has .yml extension", func(t *testing.T) {
		// given
		filePath := "values.yml"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageYAML, result)
	})

	t.Run("should return LanguageUnknown when extension is not recognized", func(t *testing.T) {
		// given
		filePath := "README.md"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageUnknown, result)
	})

	t.Run("should handle full paths with directories", func(t *testing.T) {
		// given
		filePath := "src/internal/handler.go"

		// when
		result := entities.ClassifyFileByExtension(filePath)

		// then
		assert.Equal(t, entities.LanguageGo, result)
	})
}

func TestClassifyFilesByExtension(t *testing.T) {
	t.Parallel()

	t.Run("should return unique languages for multiple files", func(t *testing.T) {
		// given
		paths := []string{"main.go", "utils.go", "index.js", "script.py"}

		// when
		result := entities.ClassifyFilesByExtension(paths)

		// then
		assert.Len(t, result, 3)
		assert.Contains(t, result, entities.LanguageGo)
		assert.Contains(t, result, entities.LanguageNode)
		assert.Contains(t, result, entities.LanguagePython)
	})

	t.Run("should exclude unknown extensions", func(t *testing.T) {
		// given
		paths := []string{"main.go", "README.md", "LICENSE"}

		// when
		result := entities.ClassifyFilesByExtension(paths)

		// then
		assert.Len(t, result, 1)
		assert.Contains(t, result, entities.LanguageGo)
	})

	t.Run("should return empty slice when no files match", func(t *testing.T) {
		// given
		paths := []string{"README.md", "LICENSE", "Makefile"}

		// when
		result := entities.ClassifyFilesByExtension(paths)

		// then
		assert.Empty(t, result)
	})

	t.Run("should deduplicate languages from multiple files with same extension", func(t *testing.T) {
		// given
		paths := []string{"main.go", "handler.go", "service.go"}

		// when
		result := entities.ClassifyFilesByExtension(paths)

		// then
		assert.Len(t, result, 1)
		assert.Contains(t, result, entities.LanguageGo)
	})
}
