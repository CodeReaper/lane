package translations

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type AndroidLanguageFile struct {
	file *LanguageFile
}

func (f *AndroidLanguageFile) Write(translations *Translations) error {
	return f.file.write(f, translations)
}

func (f *AndroidLanguageFile) write(translation *Translation, io io.Writer) error {
	regex := regexp.MustCompile(`%([0-9]+)`)

	escape := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "\\\"",
		"'", "\\'",
		"\n", "\\n")

	_, err := io.Write([]byte("<resources>\n"))
	if err != nil {
		return err
	}

	for _, k := range translation.keys {
		key := strings.ToLower(k)
		value := regex.ReplaceAllString(translation.get(k), "%${1}$$s")
		_, err = io.Write([]byte(fmt.Sprintf("\t<string name=\"%s\">%s</string>\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	_, err = io.Write([]byte("</resources>\n"))

	return err
}
