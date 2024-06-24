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

	for _, k := range translation.keys {
		key := strings.ToUpper(k)
		value := translation.get(k)
		_, err := io.Write([]byte(fmt.Sprintf("\"%s\" = \"%s\";\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	return nil
}
