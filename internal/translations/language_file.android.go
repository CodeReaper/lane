package translations

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type androidLanguageFile struct {
	file *languageFile
}

func (f *androidLanguageFile) Write(translations *translationData) error {
	return f.file.write(f, translations)
}

func (f *androidLanguageFile) write(translation *translation, io io.Writer) error {
	regex := regexp.MustCompile(`%([0-9]+)`)

	escape := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "\\\"",
		"'", "\\'",
		"\n", "\\n")

	_, err := io.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<resources>\n"))
	if err != nil {
		return err
	}

	for _, k := range translation.keys {
		key := strings.ToLower(k)
		value := regex.ReplaceAllString(translation.get(k), "%${1}$$s")
		_, err = io.Write([]byte(fmt.Sprintf("    <string name=\"%s\" formatted=\"false\">%s</string>\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	_, err = io.Write([]byte("</resources>\n"))

	return err
}
