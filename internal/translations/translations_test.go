package translations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTranslationsMissingFile(t *testing.T) {
	_, err := newTranslations("test/does-not-exist.csv")

	assert.Error(t, err)
}

func TestNewTranslations(t *testing.T) {
	translations, err := newTranslations("testdata/input.csv")

	if !assert.NoError(t, err) {
		return
	}
	assert.NotEqual(t, 0, len(translations.data))
}

func TestTranslationsTranslation(t *testing.T) {
	translations, err := newTranslations("testdata/input.csv")

	if !assert.NoError(t, err) {
		return
	}
	translation := translations.translation(1, 3)

	assert.NotNil(t, translation)
}
