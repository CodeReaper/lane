package translations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var androidKind = "android"
var iosKind = "ios"
var jsonKind = "json"

var validKinds = []string{
	iosKind,
	androidKind,
	jsonKind,
}

func Generate(ctx context.Context, flags *Flags, configurations []string) error {
	err := flags.validate()
	if err != nil {
		return err
	}

	files, err := newLanguageFiles(flags.Index, configurations)
	if err != nil {
		return err
	}

	translations, err := newTranslations(flags.Input)
	if err != nil {
		return err
	}

	kind := flags.Kind
	once := true
	for _, f := range files {
		t := translations.translation(f.keyIndex, f.valueIndex)

		switch kind {
		case androidKind:
			err = write(toAndroid, t, f.path)
		case iosKind:
			err = write(toIos, t, f.path)
			if once && err == nil {
				once = false
				err = write(toSwift, translations.translation(flags.Index, flags.DefaultIndex), flags.Output)
			}
		case jsonKind:
			err = write(toJson, t, f.path)
		default:
			err = fmt.Errorf("found unknown kind: %v", kind)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func write(writerFunc func(translation *Translation, io io.Writer) error, translation *Translation, path string) error {
	tempPath := path + ".tmp"
	defer os.Remove(tempPath)

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	err = writerFunc(translation, tempFile)
	if err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

func toAndroid(translation *Translation, io io.Writer) error {
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

func toIos(translation *Translation, io io.Writer) error {
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

func toJson(translation *Translation, io io.Writer) error {
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

func toSwift(translation *Translation, io io.Writer) error {
	regex := regexp.MustCompile(`%([0-9]+)`)

	header := `// swiftlint:disable all
import Foundation
struct Translations {
`
	footer := `}
`

	_, err := io.Write([]byte(header))
	if err != nil {
		return err
	}

	for _, k := range translation.keys {
		key := strings.ToUpper(k)
		value := translation.get(k)

		var line string
		matches := regex.FindAllStringSubmatch(value, -1)
		if len(matches) > 0 {
			arguments := make([]string, 0)
			replacements := make([]string, 0)
			for _, m := range matches {
				match := m[1]
				arguments = append(arguments, fmt.Sprintf("_ p%s: String", match))
				replacements = append(replacements, fmt.Sprintf(".replacingOccurrences(of: \"%%%s\", with: p%s)", match, match))
			}
			argumentsString := strings.Join(arguments, ", ")
			replacementsString := strings.Join(replacements, "")
			line = fmt.Sprintf("\tstatic func %s(%s) -> String { return NSLocalizedString(\"%s\", comment: \"\")%s }\n", key, argumentsString, key, replacementsString)
		} else {
			line = fmt.Sprintf("\tstatic let %s = NSLocalizedString(\"%s\", comment: \"\")\n", key, key)
		}

		_, err := io.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	_, err = io.Write([]byte(footer))
	return err
}
