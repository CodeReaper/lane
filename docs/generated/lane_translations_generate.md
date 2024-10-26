## lane translations generate

Generate translations files from a csv file

### Synopsis

Reads a CSV file and uses configuration strings to generate static resource files for android or ios.

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


```
lane translations generate [flags]
```

### Options

```
  -c, --configuration stringArray   A configuration string consisting of space separated row index and output path. Multiple configurations can be added, but one is required
  -l, --fill-in                     Fill in the value from the main/default language if a value is missing for the current language
  -h, --help                        help for generate
  -k, --index int                   The index of the key row, defaults to 0
  -i, --input string                Path to a CSV file containing a key row and a row for each language (Required)
  -m, --main-index int              Required for ios or when using fill-in. The index of the main/default language row, defaults to 0
  -o, --output string               Required for ios. A path for the generated output
  -t, --type string                 The type of output to generate, valid options are 'ios' or 'android' (Required)
```

### SEE ALSO

* [lane translations](lane_translations.md)	 - Manage translations

