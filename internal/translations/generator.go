package translations

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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

type Generator struct {
	flags *Flags
}

type Configuration struct {
	index int
	path  string
}

func NewGenerator(f *Flags) *Generator {
	return &Generator{
		flags: f,
	}
}

func (g *Generator) Generate(ctx context.Context) error {
	if err := g.flags.validate(); err != nil {
		return err
	}

	configurations, err := parseConfigurations(g.flags.Configurations)
	if err != nil {
		return err
	}

	records, err := loadRecords(g.flags.Input)
	if err != nil {
		return err
	}

	for _, c := range configurations {
		if err := writeFile(g.flags.Kind, makeTranslations(records, g.flags.Index-1, c.index-1), c.path); err != nil {
			return nil
		}
	}

	if g.flags.Kind == iosKind {
		return writeSupportFile(makeTranslations(records, g.flags.Index-1, g.flags.DefaultIndex-1), g.flags.Output)
	} else {
		return nil
	}
}

func loadRecords(path string) ([][]string, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(r)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		records = records[1:] // skips header row
	}
	return records, nil
}

func makeTranslations(records [][]string, keyIndex int, valueIndex int) map[string]string {
	translations := map[string]string{}
	for _, r := range records {
		translations[strings.ToLower(r[keyIndex])] = r[valueIndex]
	}
	return translations
}

func sortedKeys(array map[string]string) []string {
	keys := make([]string, 0, len(array))
	for k := range array {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func parseConfigurations(configurations []string) ([]Configuration, error) {
	list := make([]Configuration, 0)
	for _, configuration := range configurations {
		fields := strings.Fields(configuration)
		if len(fields) != 2 {
			return nil, fmt.Errorf("configuration has invalid format: %s", configuration)
		}

		index, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("configuration has invalid index: %s", configuration)
		}

		path := fields[1]
		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			return nil, fmt.Errorf("configuration has invalid path: %s", configuration)
		}

		list = append(list, Configuration{
			index: index,
			path:  path,
		})
	}
	return list, nil
}

func writeFile(kind string, translations map[string]string, path string) error {
	tempPath := path + ".tmp"
	defer os.Remove(tempPath)

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	switch kind {
	case androidKind:
		err = writeAndroid(translations, tempFile)
	case iosKind:
		err = writeIos(translations, tempFile)
	case jsonKind:
		err = writeJson(translations, tempFile)
	default:
		err = fmt.Errorf("found unknown kind: %v", kind)
	}
	if err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

func writeAndroid(translations map[string]string, io io.Writer) error {
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

	for _, k := range sortedKeys(translations) {
		key := strings.ToLower(k)
		value := regex.ReplaceAllString(translations[k], "%${1}$$s")
		_, err = io.Write([]byte(fmt.Sprintf("\t<string name=\"%s\">%s</string>\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	_, err = io.Write([]byte("</resources>\n"))

	return err
}

func writeIos(translations map[string]string, io io.Writer) error {
	escape := strings.NewReplacer(
		"\"", "\\\"",
		"\n", "\\n")

	for _, k := range sortedKeys(translations) {
		key := strings.ToUpper(k)
		value := translations[k]
		_, err := io.Write([]byte(fmt.Sprintf("\"%s\" = \"%s\";\n", key, escape.Replace(value))))
		if err != nil {
			return err
		}
	}

	return nil
}

func writeJson(translations map[string]string, io io.Writer) error {
	data := map[string]string{}
	for k, v := range translations {
		data[strings.ToLower(k)] = v
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = io.Write(b)

	return err
}

func writeSupportFile(translations map[string]string, path string) error {
	tempPath := path + ".tmp"
	defer os.Remove(tempPath)

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	err = writeSwift(translations, tempFile)
	if err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

func writeSwift(translations map[string]string, io io.Writer) error {
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

	for _, k := range sortedKeys(translations) {
		key := strings.ToUpper(k)
		value := translations[k]

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
