NAME
```
    shell-github-action-semver-compare
    - a lane action
```

SYNOPSIS
```
    -m main-version -c current-version [-q]
    -h
```

DESCRIPTION
```
    Compares two semver-style versions.
    The exit code will indicate if the current version is considered higher than the main version.
    The output includes a GitHub-Action-style group text for easier debugging, and an error message when exit code > 0.

    The purpose is to enable sanity testing required version changes.
```

OPTIONS
```
    -h
        Shows the full help.

    -m
        The main version

    -c
        The current version

    -q
        Quiet mode will not output debugging messages
```

EXIT CODES
```
    10
        Indicates the versions match each other.
    20
        Indicates the main version is considered higher than the current version.
```
