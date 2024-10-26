package translations

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/template"
)

type iosSupportLanguageFile struct {
	file     *languageFile
	template string
}

func (f *iosSupportLanguageFile) Write(translations *translationData) error {
	return f.file.write(f, translations)
}

type line struct {
	Name         string
	Arguments    string
	Replacements string
}

const iosSupportTemplate = `// swiftlint:disable all
import Foundation
public struct Translations {
{{- range . }}
{{- if .Arguments }}
	static func {{ .Name }}({{ .Arguments }}) -> String { return NSLocalizedString("{{ .Name }}", comment: ""){{ .Replacements }} }
{{- else }}
	static let {{ .Name }} = NSLocalizedString("{{ .Name }}", comment: "")
{{- end }}
{{- end }}
}
`

func (f *iosSupportLanguageFile) write(translation *translation, io io.Writer) error {
	regex := regexp.MustCompile(`%([0-9]+)`)

	tmpl := f.template
	if tmpl == "" {
		tmpl = iosSupportTemplate
	}

	generator, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return err
	}

	list := make([]*line, 0)
	for _, k := range translation.keys {
		key := strings.ToUpper(k)
		value := translation.get(k)
		var item *line

		matches := regex.FindAllStringSubmatch(value, -1)
		if len(matches) > 0 {
			arguments := make([]string, 0)
			replacements := make([]string, 0)
			for _, m := range matches {
				match := m[1]
				arguments = append(arguments, fmt.Sprintf("_ p%s: String", match))
				replacements = append(replacements, fmt.Sprintf(".replacingOccurrences(of: \"%%%s\", with: p%s)", match, match))
			}
			item = &line{
				Name:         key,
				Arguments:    strings.Join(arguments, ", "),
				Replacements: strings.Join(replacements, ""),
			}
		} else {
			item = &line{
				Name: key,
			}
		}

		list = append(list, item)
	}

	return generator.Execute(io, list)
}
