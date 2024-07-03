package translations

import (
	"encoding/json"
	"io"
	"strings"
)

type jsonLanguageFile struct {
	file *languageFile
}

func (f *jsonLanguageFile) Write(translations *translationData) error {
	return f.file.write(f, translations)
}

func (f *jsonLanguageFile) write(translation *translation, io io.Writer) error {
	data := map[string]string{}
	for _, k := range translation.keys {
		data[strings.ToLower(k)] = translation.get(k)
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = io.Write(b)

	return err
}
