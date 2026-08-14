//go:build unit

package dart_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dart"
)

// commentedPubspec mirrors the shape of a real Flutter manifest: rationale
// comments between the dependency entries, an SDK constraint and SDK-sourced
// packages. It exists to prove the rewrite is line-oriented — a marshal
// round-trip would silently discard every comment below.
const commentedPubspec = `name: medhub
description: 'medhub clinic management — Flutter (Material 3)'
publish_to: 'none'
version: 0.1.0

environment:
  sdk: ^3.13.0

dependencies:
  flutter:
    sdk: flutter

  # Routing. Keeps the SPA's paths verbatim and gives the authenticated-route
  # redirect that protected_route.tsx performs today.
  go_router: ^17.0.0

  intl: 0.20.2

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^6.0.0

flutter:
  uses-material-design: true
`

func writePubspec(t *testing.T, content string) (dir string, path string) {
	t.Helper()
	dir = t.TempDir()
	path = filepath.Join(dir, "pubspec.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return dir, path
}

func TestDartVersionReader(t *testing.T) {
	t.Parallel()

	t.Run("should read the version from pubspec.yaml", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, commentedPubspec)
		r := &dart.VersionReader{}

		// when
		version, err := r.ReadVersion(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, "0.1.0", version.String())
	})

	t.Run("should preserve the build number when reading", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, "name: app\nversion: 1.0.0+7\n")
		r := &dart.VersionReader{}

		// when
		version, err := r.ReadVersion(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.0.0+7", version.String())
	})

	t.Run("should return an error when the manifest declares no version", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, "name: workspace_member\nenvironment:\n  sdk: ^3.13.0\n")
		r := &dart.VersionReader{}

		// when
		_, err := r.ReadVersion(dir)

		// then
		require.Error(t, err)
	})

	t.Run("should return an error when there is no manifest", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		r := &dart.VersionReader{}

		// when
		_, err := r.ReadVersion(dir)

		// then
		require.Error(t, err)
	})

	t.Run("should report pubspec.yaml as the version file", func(t *testing.T) {
		t.Parallel()

		// given
		r := &dart.VersionReader{}

		// when
		files := r.VersionFiles()

		// then
		assert.Equal(t, []string{"pubspec.yaml"}, files)
	})
}

func TestDartVersionWriter(t *testing.T) {
	t.Parallel()

	t.Run("should update the version and preserve every comment", func(t *testing.T) {
		t.Parallel()

		// given
		dir, path := writePubspec(t, commentedPubspec)
		w := &dart.VersionWriter{}
		version, err := entities.NewVersion("0.2.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, version)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		result := string(content)
		assert.Contains(t, result, "version: 0.2.0\n")
		assert.NotContains(t, result, "version: 0.1.0")
		assert.Contains(t, result, "  # Routing. Keeps the SPA's paths verbatim and gives the authenticated-route\n")
		assert.Contains(t, result, "  sdk: ^3.13.0\n")
		assert.Contains(t, result, "  go_router: ^17.0.0\n")
		assert.Contains(t, result, "  flutter_lints: ^6.0.0\n")
	})

	t.Run("should increment the build number when the manifest carries one", func(t *testing.T) {
		t.Parallel()

		// given
		dir, path := writePubspec(t, "name: app\nversion: 1.0.0+1\n")
		w := &dart.VersionWriter{}
		version, err := entities.NewVersion("1.1.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, version)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		assert.Equal(t, "name: app\nversion: 1.1.0+2\n", string(content))
	})

	t.Run("should return an error when the manifest declares no version", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, "name: workspace_member\n")
		w := &dart.VersionWriter{}
		version, err := entities.NewVersion("1.0.0")
		require.NoError(t, err)

		// when
		err = w.WriteVersion(dir, version)

		// then
		require.Error(t, err)
	})

	t.Run("should report pubspec.yaml as the file changed", func(t *testing.T) {
		t.Parallel()

		// given
		dir, path := writePubspec(t, commentedPubspec)
		w := &dart.VersionWriter{}

		// when
		files, err := w.FilesChanged(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{path}, files)
	})

	t.Run("should return an error from FilesChanged when there is no manifest", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		w := &dart.VersionWriter{}

		// when
		_, err := w.FilesChanged(dir)

		// then
		require.Error(t, err)
	})
}
