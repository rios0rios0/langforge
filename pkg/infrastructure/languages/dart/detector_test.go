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

func TestDartDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect Dart project when pubspec.yaml exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: app\n"), 0o600))
		d := &dart.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should not detect Dart project when markers are absent", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &dart.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguageDart", func(t *testing.T) {
		t.Parallel()

		// given
		d := &dart.Detector{}

		// when
		lang := d.Language()

		// then
		assert.Equal(t, entities.LanguageDart, lang)
	})
}

func TestIsFlutterManifest(t *testing.T) {
	t.Parallel()

	t.Run("should report Flutter when the manifest has a flutter section", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\n\nflutter:\n  uses-material-design: true\n"

		// when
		result := dart.IsFlutterManifest(content)

		// then
		assert.True(t, result)
	})

	t.Run("should report Flutter when a dependency is sourced from the SDK", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\n\ndependencies:\n  flutter:\n    sdk: flutter\n"

		// when
		result := dart.IsFlutterManifest(content)

		// then
		assert.True(t, result)
	})

	t.Run("should not report Flutter for a plain Dart package", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: cli\n\ndependencies:\n  args: ^2.7.0\n  flutter_lints: ^6.0.0\n"

		// when
		result := dart.IsFlutterManifest(content)

		// then
		assert.False(t, result)
	})

	t.Run("should read the manifest from disk when given a repo path", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		content := "name: app\n\ndependencies:\n  flutter:\n    sdk: flutter\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte(content), 0o600))

		// when
		result := dart.IsFlutter(dir)

		// then
		assert.True(t, result)
	})

	t.Run("should report not Flutter when there is no manifest at all", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()

		// when
		result := dart.IsFlutter(dir)

		// then
		assert.False(t, result)
	})
}
