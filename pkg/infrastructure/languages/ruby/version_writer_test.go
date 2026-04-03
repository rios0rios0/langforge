//go:build unit

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

func TestRubyVersionWriter(t *testing.T) {
	t.Parallel()

	t.Run("should update version in .gemspec file", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		original := "Gem::Specification.new do |s|\n  s.name = \"mygem\"\n  s.version = \"1.0.0\"\nend\n"
		gemspecPath := filepath.Join(dir, "mygem.gemspec")
		require.NoError(t, os.WriteFile(gemspecPath, []byte(original), 0o600))
		w := &ruby.VersionWriter{}
		newVersion, err := entities.NewVersion("2.0.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, newVersion)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(gemspecPath)
		require.NoError(t, readErr)
		assert.Contains(t, string(content), `"2.0.0"`)
		assert.NotContains(t, string(content), `"1.0.0"`)
	})

	t.Run("should write version to VERSION file when no .gemspec exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		versionFile := filepath.Join(dir, "VERSION")
		require.NoError(t, os.WriteFile(versionFile, []byte("1.0.0\n"), 0o600))
		w := &ruby.VersionWriter{}
		newVersion, err := entities.NewVersion("3.0.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, newVersion)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(versionFile)
		require.NoError(t, readErr)
		assert.Equal(t, "3.0.0\n", string(content))
	})

	t.Run("should return error when no .gemspec or VERSION file exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		w := &ruby.VersionWriter{}
		newVersion, err := entities.NewVersion("1.0.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, newVersion)

		// then
		require.Error(t, err)
	})

	t.Run("should return correct files changed for .gemspec", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		gemspecPath := filepath.Join(dir, "mygem.gemspec")
		require.NoError(t, os.WriteFile(gemspecPath, []byte("Gem::Specification.new {}\n"), 0o600))
		w := &ruby.VersionWriter{}

		// when
		files, err := w.FilesChanged(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{gemspecPath}, files)
	})

	t.Run("should return VERSION in files changed when no .gemspec exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		versionFile := filepath.Join(dir, "VERSION")
		require.NoError(t, os.WriteFile(versionFile, []byte("1.0.0\n"), 0o600))
		w := &ruby.VersionWriter{}

		// when
		files, err := w.FilesChanged(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{versionFile}, files)
	})
}
