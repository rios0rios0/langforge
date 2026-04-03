//go:build unit

package ruby_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/ruby"
)

func TestRubyVersionReader(t *testing.T) {
	t.Parallel()

	t.Run("should read version from .gemspec file", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		content := "Gem::Specification.new do |s|\n  s.name = \"mygem\"\n  s.version = \"1.2.3\"\nend\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "mygem.gemspec"), []byte(content), 0o600))
		r := &ruby.VersionReader{}

		// when
		version, err := r.ReadVersion(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", version.String())
	})

	t.Run("should read version from VERSION file when no .gemspec exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "VERSION"), []byte("2.0.0\n"), 0o600))
		r := &ruby.VersionReader{}

		// when
		version, err := r.ReadVersion(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, "2.0.0", version.String())
	})

	t.Run("should return error when no version found", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		r := &ruby.VersionReader{}

		// when
		_, err := r.ReadVersion(dir)

		// then
		require.Error(t, err)
	})

	t.Run("should return error when directory does not exist", func(t *testing.T) {
		t.Parallel()

		// given
		r := &ruby.VersionReader{}

		// when
		_, err := r.ReadVersion("/nonexistent/path")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reading directory")
	})
}
