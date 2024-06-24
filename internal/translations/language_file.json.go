package translations

import (
	"encoding/json"
	"io"
	"strings"
)

type JsonLanguageFile struct {
	file *LanguageFile
}

func (f *JsonLanguageFile) Write(translations *Translations) error {
	return f.file.write(f, translations)
}

func (f *JsonLanguageFile) write(translation *Translation, io io.Writer) error {
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
