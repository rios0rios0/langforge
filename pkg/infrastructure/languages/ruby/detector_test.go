package ruby_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/ruby"
)

func TestRubyDetector(t *testing.T) {
	t.Parallel()

	markers := []struct {
		name    string
		file    string
		content string
	}{
		{
			name:    "should detect Ruby project when Gemfile exists",
			file:    "Gemfile",
			content: "source 'https://rubygems.org'\n",
		},
		{
			name:    "should detect Ruby project when .gemspec exists",
			file:    "mygem.gemspec",
			content: "Gem::Specification.new {}\n",
		},
	}
	for _, marker := range markers {
		t.Run(marker.name, func(t *testing.T) {
			t.Parallel()

			// given
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, marker.file), []byte(marker.content), 0o600))
			d := &ruby.Detector{}

			// when
			detected, err := d.Detect(dir)

			// then
			require.NoError(t, err)
			assert.True(t, detected)
		})
	}

	t.Run("should not detect Ruby project when markers are absent", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &ruby.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguageRuby", func(t *testing.T) {
		t.Parallel()

		// given
		d := &ruby.Detector{}

		// when
		lang := d.Language()

		// then
		assert.Equal(t, entities.LanguageRuby, lang)
	})
}
