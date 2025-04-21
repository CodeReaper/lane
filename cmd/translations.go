package cmd

import (
	"context"

	"github.com/codereaper/lane/internal/downloader"
	"github.com/codereaper/lane/internal/translations"
	"github.com/spf13/cobra"
)

func newTranslationsCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "translations",
		Short: "Manage translations",
		Long:  "Download translations from google sheets and/or generate translations from local csv files.",
	}
	cmd.AddCommand(newTranslationsDownloadCommand())
	cmd.AddCommand(newTranslationsGenerateCommand())
	return cmd
}

func newTranslationsDownloadCommand() *cobra.Command {
	var additionalHelp = `Authentication is done using a json file issued by Google. You get this json file by creating a "Service Account Key", which if you do not have a service account, requires you to first create a service account.

Creating both an account and a key is explaining here: https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount

You may have to enable Google Drive API access when using it for the first time. The error message(s) should provide a direct link to enabling access.

Make sure to share the sheet with the 'client_email' assigned to your service account.
`
	var flags downloader.Flags
	var cmd = &cobra.Command{
		Use:     "download",
		Short:   "Download translations",
		Long:    additionalHelp,
		Example: "  lane translations download -o output.csv -c google-api.json -d 11p...ev7lc -f csv",
		RunE: func(cmd *cobra.Command, args []string) error {
			return downloader.Download(context.Background(), &flags)
		},
	}
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Path to save output file (Required)")
	cmd.Flags().StringVarP(&flags.Credentials, "credentials", "c", "", "A path to the credentials json file issued by Google (Required). More details under help")
	cmd.Flags().StringVarP(&flags.DocumentId, "document", "d", "", "The document id of the sheet to download (Required). Found in its url, e.g. https://docs.google.com/spreadsheets/d/<document-id>/edit#gid=0")
	cmd.Flags().StringVarP(&flags.Format, "format", "f", "", "The format of the output, defaults to csv")
	cmd.MarkFlagRequired("output")
	cmd.MarkFlagRequired("credentials")
	cmd.MarkFlagRequired("document")
	return cmd
}

func newTranslationsGenerateCommand() *cobra.Command {
	var additionalHelp = `Reads a CSV file and uses configuration strings to generate static resource files for android and ios.

JSON output is also supported as simple output, the placeholders discussed below do not apply to JSON output.

Each translated string can have '%<digit>'-style placeholders, however the number of placeholder for each translated language must be the same.
The placeholders in the generated output will always take a string as input.

The purpose is to enable compilation checks for translated strings with an external source for the actual strings.

EXAMPLES:

If the contents of 'input.csv' is:

    KEY,UPDATE NEEDED,English,Danish,COMMENT
    SOMETHING,,Something,Noget,
    SOMETHING_WITH_ARGUMENTS,,Something with %1 and %2,Noget med %1 og %2,

- Android

The output using '-t android -i input.csv -c "3 en.xml" -k 1' would be:

    <resources>
            <string name="something">Something</string>
            <string name="something_with_arguments">Something with %1$s and %2$s</string>
    </resources>

- iOS

The output using '-t ios -i input.csv -c "3 en.strings" -k 1 -m 3 -o translations.swift' would be:

en.strings:

    "SOMETHING" = "Something";
    "SOMETHING_WITH_ARGUMENTS" = "Something with %1 and %2";

translations.swift:

    // swiftlint:disable all
    import Foundation
    struct Translations {
            static let SOMETHING = NSLocalizedString("SOMETHING", comment: "")
            static func SOMETHING_WITH_ARGUMENTS(_ p1: String, _ p2: String) -> String { return NSLocalizedString("SOMETHING_WITH_ARGUMENTS", comment: "").replacingOccurrences(of: "%1", with: p1).replacingOccurrences(of: "%2", with: p2) }
    }

You can support your own golang text template to change the output, the above output is generated with the following default template:

	// swiftlint:disable all
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
	var flags translations.Flags
	var configurations []string
	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate translations files from a csv file",
		Long:  additionalHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return translations.Generate(context.Background(), &flags, configurations)
		},
	}
	cmd.Flags().StringVarP(&flags.Input, "input", "i", "", "Path to a CSV file containing a key row and a row for each language (Required)")
	cmd.Flags().StringVarP(&flags.Kind, "type", "t", "", "The type of output to generate, valid options are 'ios', 'android' or 'json' (Required)")
	cmd.Flags().IntVarP(&flags.KeyIndex, "index", "k", 0, "The index of the key row (Required)")
	cmd.Flags().StringArrayVarP(&configurations, "configuration", "c", make([]string, 0), "A configuration string consisting of space separated row index and output path. Multiple configurations can be added, but one is required")
	cmd.Flags().IntVarP(&flags.DefaultValueIndex, "main-index", "m", 0, "Required by type ios and by option fill-in. The index of the main/default language row")
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Required for type ios. A path for the generated output")
	cmd.Flags().StringVarP(&flags.Template, "template", "p", "", "Only for type ios and optional. A path for the template to generate from")
	cmd.Flags().BoolVarP(&flags.FillIn, "fill-in", "l", false, "Fill in the value from the main/default language if a value is missing for the current language")
	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("configuration")
	return cmd
}
