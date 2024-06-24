package translations

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIosLanguageFileInterface(t *testing.T) {
	var f interface{} = iosLanguageFile{}
	w, ok := f.(languageFileWriter)
	assert.True(t, ok)
	assert.NotNil(t, w)
}

func TestIosLanguageFileWriteEmpty(t *testing.T) {
	tr := newTranslation(map[string]string{})
	f := &iosLanguageFile{}
	var b bytes.Buffer

	err := f.write(tr, &b)
	assert.NoError(t, err)
}

func TestIosLanguageFileWriteInputFile(t *testing.T) {
	translations, err := newTranslations("testdata/input.csv")
	if !assert.NoError(t, err) {
		return
	}

	tr := translations.translation(1, 3)

	f := &iosLanguageFile{}
	var b bytes.Buffer

	err = f.write(tr, &b)
	assert.NoError(t, err)
	t.Fail()
}
