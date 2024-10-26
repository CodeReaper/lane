package translations

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var interfaces = []interface{}{
	&androidLanguageFile{},
	&iosSupportLanguageFile{},
	&iosLanguageFile{},
	&jsonLanguageFile{},
}

func TestLanguageFileInterface(t *testing.T) {
	for _, x := range interfaces {
		w, ok := x.(languageFileWriter)
		assert.True(t, ok)
		assert.NotNil(t, w)
	}
}

var emptyWriters = []struct {
	writer       languageFileWriter
	expectedPath string
}{
	{
		&androidLanguageFile{},
		"testdata/android-empty.expected",
	},
	{
		&iosSupportLanguageFile{},
		"testdata/ios-swift-empty.expected",
	},
	{
		&jsonLanguageFile{},
		"testdata/json-empty.expected",
	},
}

func TestLanguageFileWriteEmpty(t *testing.T) {
	for _, x := range emptyWriters {
		tr := newTranslation(map[string]string{})
		var b bytes.Buffer

		err := x.writer.write(tr, &b)
		assert.NoError(t, err)

		expected, err := os.ReadFile(x.expectedPath)
		if !assert.NoError(t, err) {
			return
		}
		assert.EqualValues(t, expected, b.Bytes())
	}
}

func TestIosLanguageFileWriteEmpty(t *testing.T) {
	x := &iosLanguageFile{}
	tr := newTranslation(map[string]string{})
	var b bytes.Buffer

	err := x.write(tr, &b)
	assert.NoError(t, err)
	assert.Nil(t, b.Bytes())
}

var inputWriters = []struct {
	writer       languageFileWriter
	index        int
	expectedPath string
}{
	{
		&androidLanguageFile{},
		3,
		"testdata/android-en.expected",
	},
	{
		&androidLanguageFile{},
		4,
		"testdata/android-da.expected",
	},
	{
		&iosLanguageFile{},
		3,
		"testdata/ios-en.expected",
	},
	{
		&iosLanguageFile{},
		4,
		"testdata/ios-da.expected",
	},
	{
		&iosSupportLanguageFile{},
		3,
		"testdata/ios-swift.expected",
	},
	{
		&jsonLanguageFile{},
		3,
		"testdata/json-en.expected",
	},
}

func TestLanguageFileWriteInputFile(t *testing.T) {
	translations, err := newTranslations("testdata/input.csv")
	if !assert.NoError(t, err) {
		return
	}

	for _, x := range inputWriters {
		tr := translations.translation(1, x.index, false, 0)
		var b bytes.Buffer

		err := x.writer.write(tr, &b)
		assert.NoError(t, err)

		expected, err := os.ReadFile(x.expectedPath)
		if !assert.NoError(t, err) {
			return
		}
		assert.EqualValues(t, expected, b.Bytes())
	}
}

func TestEscapingSingleColumn(t *testing.T) {
	output := "<resources>\n\t<string name=\"that\">Another string including a | even.</string>\n\t<string name=\"this\">This is a longer sentence, which includes a comma.</string>\n</resources>\n"
	outputPath := "../../build/out.csv"
	inputPath := "../../build/test.csv"
	err := os.WriteFile("../../build/test.csv", []byte("key,en\nTHIS,\"This is a longer sentence, which includes a comma.\"\nTHAT,Another string including a | even."), 0777)
	if !assert.NoError(t, err) {
		return
	}

	flags := Flags{
		Input:    inputPath,
		Kind:     "android",
		KeyIndex: 1,
	}
	configurations := []string{"2 " + outputPath}
	err = Generate(context.Background(), &flags, configurations)
	if !assert.NoError(t, err) {
		return
	}

	b, err := os.ReadFile(outputPath)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, output, string(b))
}

func TestEscapingAllColumns(t *testing.T) {
	output := "<resources>\n\t<string name=\"that\">Another string including a | even.</string>\n\t<string name=\"this\">This is a longer sentence, which includes a comma.</string>\n</resources>\n"
	outputPath := "../../build/out.csv"
	inputPath := "../../build/test.csv"
	err := os.WriteFile("../../build/test.csv", []byte("\"key\",\"en\"\n\"THIS\",\"This is a longer sentence, which includes a comma.\"\n\"THAT\",\"Another string including a | even.\""), 0777)
	if !assert.NoError(t, err) {
		return
	}

	flags := Flags{
		Input:    inputPath,
		Kind:     "android",
		KeyIndex: 1,
	}
	configurations := []string{"2 " + outputPath}
	err = Generate(context.Background(), &flags, configurations)
	if !assert.NoError(t, err) {
		return
	}

	b, err := os.ReadFile(outputPath)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, output, string(b))
}
