package translations

import (
	"fmt"
	"io"
	"strings"
)

type iosLanguageFile struct {
	file *languageFile
}

func (f *iosLanguageFile) Write(translations *translationData) error {
	return f.file.write(f, translations)
}

func (f *iosLanguageFile) write(translation *translation, io io.Writer) error {
	escape := strings.NewReplacer(
		"\"", "\\\"",
		"\n", "\\n")

	for _, key := range translation.keys {
		value := translation.get(key)
		_, err := io.Write([]byte(fmt.Sprintf("\"%s\" = \"%s\";\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	return nil
}
