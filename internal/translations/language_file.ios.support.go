package translations

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type IosSupportLanguageFile struct {
	file *LanguageFile
}

func (f *IosSupportLanguageFile) Write(translations *Translations) error {
	return f.file.write(f, translations)
}

func (f *IosSupportLanguageFile) write(translation *Translation, io io.Writer) error {
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
