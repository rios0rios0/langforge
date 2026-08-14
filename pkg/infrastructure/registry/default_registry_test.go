//go:build unit

package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/registry"
)

func TestNewDefaultRegistry(t *testing.T) {
	t.Parallel()

	t.Run("should detect a Dart project by its pubspec.yaml", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: app\n"), 0o600))
		reg := registry.NewDefaultRegistry()

		// when
		provider, err := reg.Detect(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, entities.LanguageDart, provider.Language())
	})

	t.Run("should prefer Dart over Node when a Flutter web project also has a package.json",
		func(t *testing.T) {
			t.Parallel()

			// given — a Flutter web project that keeps a package.json for its
			// JavaScript tooling. package.json is a weak marker and pubspec.yaml
			// is not, so detection must not depend on registration luck.
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: app\n"), 0o600))
			require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}\n"), 0o600))
			reg := registry.NewDefaultRegistry()

			// when
			provider, err := reg.Detect(dir)

			// then
			require.NoError(t, err)
			assert.Equal(t, entities.LanguageDart, provider.Language())
		})

	t.Run("should still detect Node when only a package.json is present", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}\n"), 0o600))
		reg := registry.NewDefaultRegistry()

		// when
		provider, err := reg.Detect(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, entities.LanguageNode, provider.Language())
	})

	t.Run("should expose a provider for every registered language", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewDefaultRegistry()

		// when
		provider, err := reg.Get(entities.LanguageDart)

		// then
		require.NoError(t, err)
		assert.Equal(t, entities.LanguageDart, provider.Language())
	})
}
