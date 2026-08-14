//go:build unit

package dart_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dart"
)

func constraintOf(deps []entities.Dependency, name string) (string, bool) {
	for _, d := range deps {
		if d.Name == name {
			return d.Current, true
		}
	}
	return "", false
}

func TestDartDependencyReader(t *testing.T) {
	t.Parallel()

	t.Run("should read dependencies and dev dependencies", func(t *testing.T) {
		t.Parallel()

		// given
		content := commentedPubspec

		// when
		deps := dart.ParseDependencies(content)

		// then
		goRouter, ok := constraintOf(deps, "go_router")
		require.True(t, ok)
		assert.Equal(t, "^17.0.0", goRouter)

		intl, ok := constraintOf(deps, "intl")
		require.True(t, ok)
		assert.Equal(t, "0.20.2", intl)

		lints, ok := constraintOf(deps, "flutter_lints")
		require.True(t, ok)
		assert.Equal(t, "^6.0.0", lints)
	})

	t.Run("should skip packages sourced from the SDK", func(t *testing.T) {
		t.Parallel()

		// given
		content := commentedPubspec

		// when
		deps := dart.ParseDependencies(content)

		// then
		_, hasFlutter := constraintOf(deps, "flutter")
		assert.False(t, hasFlutter)
		_, hasFlutterTest := constraintOf(deps, "flutter_test")
		assert.False(t, hasFlutterTest)
	})

	t.Run("should skip packages sourced from git or a path", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\n\ndependencies:\n" +
			"  forked_pkg:\n    git:\n      url: https://example.com/forked_pkg.git\n" +
			"  local_pkg:\n    path: ../local_pkg\n" +
			"  args: ^2.7.0\n"

		// when
		deps := dart.ParseDependencies(content)

		// then
		require.Len(t, deps, 1)
		assert.Equal(t, "args", deps[0].Name)
		assert.Equal(t, "^2.7.0", deps[0].Current)
	})

	t.Run("should strip quotes and trailing comments from constraints", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\n\ndependencies:\n" +
			"  intl: '0.20.2' # pinned, see docs/i18n.md\n" +
			"  meta: \"^1.18.3\"\n"

		// when
		deps := dart.ParseDependencies(content)

		// then
		intl, ok := constraintOf(deps, "intl")
		require.True(t, ok)
		assert.Equal(t, "0.20.2", intl)

		meta, ok := constraintOf(deps, "meta")
		require.True(t, ok)
		assert.Equal(t, "^1.18.3", meta)
	})

	t.Run("should ignore keys outside the dependency sections", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\nversion: 1.0.0\n\n" +
			"environment:\n  sdk: ^3.13.0\n\n" +
			"dependencies:\n  args: ^2.7.0\n\n" +
			"flutter:\n  uses-material-design: true\n"

		// when
		deps := dart.ParseDependencies(content)

		// then
		require.Len(t, deps, 1)
		assert.Equal(t, "args", deps[0].Name)
	})

	t.Run("should record dependency overrides", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\n\ndependency_overrides:\n  collection: 1.19.1\n"

		// when
		deps := dart.ParseDependencies(content)

		// then
		require.Len(t, deps, 1)
		assert.Equal(t, "collection", deps[0].Name)
		assert.Equal(t, "1.19.1", deps[0].Current)
		assert.Equal(t, "pubspec.yaml", deps[0].SourceFile)
	})

	t.Run("should return an error when there is no manifest", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		r := &dart.DependencyReader{}

		// when
		_, err := r.ReadDependencies(dir)

		// then
		require.Error(t, err)
	})

	t.Run("should read dependencies from disk", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, commentedPubspec)
		r := &dart.DependencyReader{}

		// when
		deps, err := r.ReadDependencies(dir)

		// then
		require.NoError(t, err)
		_, ok := constraintOf(deps, "go_router")
		assert.True(t, ok)
	})
}
