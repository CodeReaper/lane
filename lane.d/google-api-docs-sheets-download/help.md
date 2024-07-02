NAME
```
    google-api-docs-sheets-download
    - a lane action
```

SYNOPSIS
```
    -o directory -i string -s string [-f string]
    -h
```

DESCRIPTION
```
    DEPRECATED - in version 1.0.0 this action is available as 'lane translations download [OPTIONS]', but will require a JSON type service account key.

    Downloads a Sheet from Google Docs.
```

OPTIONS
```
    -h
        Shows the full help.

    -o
        A path to write the output to.

    -i
        The document id of the sheet to download. Found in its url, e.g. https://docs.google.com/spreadsheets/d/<document-id>/edit#gid=0

    -t
        A JWT, see `lane google-api-jwt-generate -h`.

    -f
        The format of the output, defaults to csv.
```
