## lane translations download

Download translations

### Synopsis

Authentication is done using a json file issued by Google. You get this json file by creating a "Service Account Key", which if you do not have a service account, requires you to first create a service account.

Creating both an account and a key is explaining here: https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount

You may have to enable Google Drive API access when using it for the first time. The error message(s) should provide a direct link to enabling access.

Make sure to share the sheet with the 'client_email' assigned to your service account.


```
lane translations download [flags]
```

### Examples

```
  lane translations download -o output.csv -c google-api.json -d 11p...ev7lc -f csv
```

### Options

```
  -c, --credentials string   A path to the credentials json file issued by Google (Required). More details under help
  -d, --document string      The document id of the sheet to download (Required). Found in its url, e.g. https://docs.google.com/spreadsheets/d/<document-id>/edit#gid=0
  -f, --format string        The format of the output, defaults to csv
  -h, --help                 help for download
  -o, --output string        Path to save output file (Required)
```

### SEE ALSO

* [lane translations](lane_translations.md)	 - Manage translations

