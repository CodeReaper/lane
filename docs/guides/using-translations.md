## Using Translations

The purpose of translations in lane is to automate handling of text for Android, iOS and web applications.

Translation in lane provides these functions:

1. An optional helper function for downloading CSV files from Google Sheets.
1. A function that writes out the actual files and support files needed for each application.

### A Practical Example

Suppose you have a Google Sheet document with an ID of `MY-DOCUMENT` which contains translations for two languages, English and Danish. This documents should have contents similar to this:

| Key              | Context                 | EN       | DA     |
| ---------------- | ----------------------- | -------- | ------ |
| WELCOME_MESSAGE  | Message upon login      | Hello %1 | Hej %1 |
| BUY              | The purchase now button | Acquire  | Køb    |

Note that the string key is in column `1` and that English is in column `3` and Danish in column `4` and that the header row is only for us humans.

We will assume you already have your credentials in `lane/translations.json` (for help with credentials see [](../../docs/generated/lane_translations_download.md)).

What we will be doing to generate translation files for Android, iOS and web is to use `docker run` to invoke lane. Below are examples of how such calls should be integrated into tooling your project may already be using for easier copy-and-paste, but ...

**Remember**:
1. Update the document id, if used
1. Remove download step if you are using a local CSV file.
1. Adjust the configurations and other options to match your needs.

Alright, go ahead and copy!

#### Android

We will update the `build.gradle` file so we can update all translations by running a gradle task, or simply:

```
./gradle updateStrings
```

The updates to a `build.gradle` file should look like this:

```gradle
static void executeWithErrorHandling(ArrayList<String> command, File dir) {
    def process = command.execute(null, dir)
    process.waitFor()
    if (process.exitValue() != 0) throw new GradleException("Gradle exited with: " + process.exitValue() + "\nDetails: " + process.errorReader().text)
}

static void dockerRun(ArrayList<String> args, File dir) {
    executeWithErrorHandling(['docker', 'run', '-it', '--rm', '--quiet-pull', '-v .:/workspace', '-w /workspace', 'ghcr.io/codereaper/lane@1', '--'] + args, dir)
}

tasks.register('updateStrings') {
    dockerRun(['lane', 'translations', 'download',
            '--credentials', 'lane/translations.json', '--output', 'lane/translations.csv', '--document', 'MY-DOCUMENT'],
        projectDir)

    dockerRun(['lane', 'translations', 'generate',
            '--fill-in', '--input', 'lane/translations.csv', '--type', 'android',
            '--index', '1', '--main-index', '3',
            '--configuration', '3 src/main/res/values/strings.xml',
            '--configuration', '3 src/main/res/values-en-rGB/strings.xml',
            '--configuration', '4 src/main/res/values-da/strings.xml'],
        projectDir)
}
```

#### iOS

We will add a target to a Makefile so we can update all translations by running a single task:

```
make update-strings
```

The updates to a `Makefile` file should look like this:

```make
LANE_RUN = docker run -it --rm --quiet-pull -v .:/workspace -w /workspace ghcr.io/codereaper/lane@1 -- lane

update-strings:
    $(LANE_RUN) translations download --credentials lane/translations.json --output lane/translations.csv \
        --document MY-DOCUMENT

    $(LANE_RUN) translations generate --fill-in --input lane/translations.csv --type ios \
        --output iOSProject/Assets/Translations.swift \
        --index 1 --main-index 3 \
        --configuration, '3 iOSProject/Assets/Translations/Base.lproj/Localizable.strings' \
        --configuration, '3 iOSProject/Assets/Translations/Base.lproj/InfoPlist.strings' \
        --configuration, '3 iOSProject/Assets/Translations/en-GB.lproj/Localizable.strings' \
        --configuration, '3 iOSProject/Assets/Translations/en-GB.lproj/InfoPlist.strings' \
        --configuration, '4 iOSProject/Assets/Translations/da.lproj/Localizable.strings' \
        --configuration, '4 iOSProject/Assets/Translations/da.lproj/InfoPlist.strings'
```

#### Web

We will add a target to a Makefile so we can update all translations by running a single task:

```
make update-strings
```

The updates to a `Makefile` file should look like this:

```make
LANE_RUN = docker run -it --rm --quiet-pull -v .:/workspace -w /workspace ghcr.io/codereaper/lane@1 -- lane

update-strings:
    $(LANE_RUN) translations download --credentials lane/translations.json --output lane/translations.csv \
        --document MY-DOCUMENT

    $(LANE_RUN) translations generate --fill-in --input lane/translations.csv --type json \
        --index 1 --main-index 3 \
        --configuration, '3 src/locale/messages/en.json' \
        --configuration, '4 src/locale/messages/da.json'
```
