package translations

import (
	"fmt"
	"io"
	"strings"
)

type IosLanguageFile struct {
	file *LanguageFile
}

func (f *IosLanguageFile) Write(translations *Translations) error {
	return f.file.write(f, translations)
}

func (f *IosLanguageFile) write(translation *Translation, io io.Writer) error {
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
