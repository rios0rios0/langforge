package registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/registry"
	"github.com/rios0rios0/langforge/test/domain/entitybuilders"
)

func TestNewLanguageRegistry(t *testing.T) {
	t.Parallel()

	t.Run("should create empty registry", func(t *testing.T) {
		t.Parallel()

		// given / when
		reg := registry.NewLanguageRegistry()

		// then
		require.NotNil(t, reg)
		assert.Empty(t, reg.Languages())
	})
}

func TestLanguageRegistryRegisterAndGet(t *testing.T) {
	t.Parallel()

	t.Run("should return provider when language is registered", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewLanguageRegistry()
		stub := entitybuilders.NewLanguageProviderStubBuilder().
			WithLanguage(entities.LanguageGo).
			Build()
		reg.Register(stub)

		// when
		provider, err := reg.Get(entities.LanguageGo)

		// then
		require.NoError(t, err)
		assert.Equal(t, entities.LanguageGo, provider.Language())
	})

	t.Run("should return error when language is not registered", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewLanguageRegistry()

		// when
		_, err := reg.Get(entities.LanguageGo)

		// then
		require.Error(t, err)
	})
}

func TestLanguageRegistryDetect(t *testing.T) {
	t.Parallel()

	t.Run("should return provider when detection matches", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewLanguageRegistry()
		stub := entitybuilders.NewLanguageProviderStubBuilder().
			WithLanguage(entities.LanguageGo).
			WithDetectResult(true, nil).
			Build()
		reg.Register(stub)

		// when
		provider, err := reg.Detect("/some/path")

		// then
		require.NoError(t, err)
		assert.Equal(t, entities.LanguageGo, provider.Language())
	})

	t.Run("should return error when no language is detected", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewLanguageRegistry()
		stub := entitybuilders.NewLanguageProviderStubBuilder().
			WithLanguage(entities.LanguageGo).
			WithDetectResult(false, nil).
			Build()
		reg.Register(stub)

		// when
		_, err := reg.Detect("/some/path")

		// then
		require.Error(t, err)
	})
}

func TestLanguageRegistryLanguages(t *testing.T) {
	t.Parallel()

	t.Run("should return all registered languages", func(t *testing.T) {
		t.Parallel()

		// given
		reg := registry.NewLanguageRegistry()
		for _, lang := range []entities.Language{entities.LanguageGo, entities.LanguageNode} {
			stub := entitybuilders.NewLanguageProviderStubBuilder().
				WithLanguage(lang).
				Build()
			reg.Register(stub)
		}

		// when
		langs := reg.Languages()

		// then
		assert.Len(t, langs, 2)
		assert.ElementsMatch(t, []entities.Language{entities.LanguageGo, entities.LanguageNode}, langs)
	})
}
