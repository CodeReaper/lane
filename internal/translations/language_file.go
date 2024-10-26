package translations

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type languageFile struct {
	path          string
	keyIndex      int
	valueIndex    int
	useFallback   bool
	fallbackIndex int
}

type languageFileWriter interface {
	write(translation *translation, io io.Writer) error
	Write(translations *translationData) error
}

func (f *languageFile) write(writer languageFileWriter, translations *translationData) error {
	translation := translations.translation(f.keyIndex, f.valueIndex, f.useFallback, f.fallbackIndex)

	tempPath := f.path + ".tmp"
	defer os.Remove(tempPath)

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	err = writer.write(translation, tempFile)
	if err != nil {
		return err
	}

	return os.Rename(tempPath, f.path)
}

func newLanguageFiles(flags *Flags, configurations []string) ([]languageFileWriter, error) {
	if len(configurations) == 0 {
		return nil, fmt.Errorf("no configurations provided")
	}

	arguments := make(map[string]int, 0)
	for _, configuration := range configurations {
		fields := strings.Fields(configuration)
		if len(fields) != 2 {
			return nil, fmt.Errorf("configuration has invalid format: %s", configuration)
		}

		index, err := strconv.Atoi(fields[0])
		if err != nil || index <= 0 {
			return nil, fmt.Errorf("configuration has invalid index: %s", configuration)
		}

		path := fields[1]
		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			return nil, fmt.Errorf("configuration has invalid path: %s", configuration)
		}

		arguments[path] = index
	}

	list := make([]languageFileWriter, 0)
	for path, index := range arguments {
		var writer languageFileWriter

		file := &languageFile{
			path:          path,
			keyIndex:      flags.KeyIndex,
			valueIndex:    index,
			useFallback:   flags.FillIn,
			fallbackIndex: flags.DefaultValueIndex,
		}

		switch flags.Kind {
		case androidKind:
			writer = &androidLanguageFile{file: file}
		case iosKind:
			writer = &iosLanguageFile{file: file}
		case jsonKind:
			writer = &jsonLanguageFile{file: file}
		default:
			return nil, fmt.Errorf("found unknown kind: %v", flags.Kind)
		}

		list = append(list, writer)
	}

	switch flags.Kind {
	case iosKind:
		var template string

		if len(flags.Template) > 0 {
			buf, err := os.ReadFile(flags.Template)
			if err != nil {
				return nil, err
			}
			template = string(buf)
		}

		supporter := &iosSupportLanguageFile{
			file: &languageFile{
				path:        flags.Output,
				keyIndex:    flags.KeyIndex,
				valueIndex:  flags.DefaultValueIndex,
				useFallback: false,
			},
			template: template,
		}
		list = append(list, supporter)
	}

	return list, nil
}
