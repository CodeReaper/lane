NAME
```
    google-api-jwt-generate
    - a lane action
```

SYNOPSIS
```
    -i string -s string -p file
    -h
```

DESCRIPTION
```
    DEPRECATED - this functionality will be built into 'lane translations download [OPTIONS]' in version 1.0.0.

    Constructs JWT generation request for Google APIs and outputs a JWT.

    The purpose is to faciliate token generation for other usages of the Google APIs.
```

OPTIONS
```
    -h
        Shows the full help.

    -i
        The issuer of the JWT, e.g. <name>@<organisation>.iam.gserviceaccount.com

    -s
        The scope(s) applied to the JWT. Apply multiple scopes by space separating them.

    -p
        A path to the .p12 file issued by Google. See authentication section.
```

AUTHENTICATION

Authentication happens using a specially created .p12 file issued by Google which must match the issuer.

You get this .p12 key by creating a "Service Account Key", which if you do not have a service account, requires you to first create a service account.

Creating both an account and a key is explained here: https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount
