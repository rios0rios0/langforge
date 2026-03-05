//go:build unit

package python_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/python"
)

func TestPythonDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect Python project when pyproject.toml exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[project]\nname=\"foo\"\nversion=\"1.0.0\"\n"), 0o600))
		d := &python.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should detect Python project when setup.py exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "setup.py"), []byte("from setuptools import setup\n"), 0o600))
		d := &python.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should return LanguagePython", func(t *testing.T) {
		t.Parallel()

		d := &python.Detector{}
		assert.Equal(t, entities.LanguagePython, d.Language())
	})
}
