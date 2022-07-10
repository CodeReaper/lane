NAME
```
    mobile-update-translations
    - a lane action
```

SYNOPSIS
```
    -t android -i file -c configuration -k index
    -t ios -i file -c configuration -k index -m index -o file
    -h
```

DESCRIPTION
```
    Reads a CSV file and uses configuration strings to generate static resource files for android or ios.

    Each translated string can have '%<digit>'-style placeholders, however the number of placeholder for each translated language must be the same.
    The placeholders in the generated output will always take a string as input.

    The purpose is to enable compilation checks for translated strings with an external source for the actual strings.
```

OPTIONS
```
    -h
        Shows the full help.

    -t
        The type of output to generate, valid options are 'ios' or 'android'.

    -i
        A CSV file containing a key row and a row for each language.

    -s
        The separator used in the CSV file. Defaults to `,`. This option is considered experimental.

    -k
        The index of the key row.

    -c
        A configuration string consisting of space separated row index and output path. Multiple configurations can be added.

    -m
        Relevant for ios only. The index of the main/default language row.

    -o
        Relevant for ios only. A path for the generated output.
```

EXAMPLES

If the contents of `input.csv` is:

```csv
    KEY,UPDATE NEEDED,English,Danish,COMMENT
    SOMETHING,,Something,Noget,
    SOMETHING_WITH_ARGUMENTS,,Something with %1 and %2,Noget med %1 og %2,
```

Android
---

The output using `-t android -i input.csv -c '3 en.xml' -k 1` would be:

```xml
    <resources>
            <string name="something">Something</string>
            <string name="something_with_arguments">Something with %1$s and %2$s</string>
    </resources>
```

iOS
---

The output using `-t ios -i input.csv -c '3 en.strings' -k 1 -m 3 -o translations.swift` would be:

en.strings:

```ini
    "SOMETHING" = "Something";
    "SOMETHING_WITH_ARGUMENTS" = "Something with %1 and %2";
```

translations.swift:

```swift
    // swiftlint:disable all
    import Foundation
    struct Translations {
            static let SOMETHING = NSLocalizedString("SOMETHING", comment: "")
            static func SOMETHING_WITH_ARGUMENTS(_ p1: String, _ p2: String) -> String { return NSLocalizedString("SOMETHING_WITH_ARGUMENTS", comment: "").replacingOccurrences(of: "%1", with: p1).replacingOccurrences(of: "%2", with: p2) }
    }
```
