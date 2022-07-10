#!/bin/sh

for command in curl base64 openssl jq; do
  if ! [ -x "$(command -v $command)" ]; then
    echo "Error: $command is not installed." >&2
    exit 3
  fi
done

unset -v issuer scopes p12
while getopts "i:s:p:" option; do
  case $option in
  i) issuer="$OPTARG" ;;
  s) scopes="$OPTARG" ;;
  p) p12="$OPTARG" ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$issuer" ] || [ -z "$scopes" ] || [ -z "$p12" ]; then
  echo "Must provide issuer, scopes and an p12 file."
  exit 111
fi

if [ ! -f "$p12" ]; then
  echo "Must provide a p12 file."
  exit 4
fi

DIR=$(mktemp -dq)

trap 'set +x; rm -fr $DIR >/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

encode() {
  base64 - | tr -d '\n=' | tr '/+' '_-'
}

iat=$(date +"%s")
exp=$(date -d +5mins +"%s")
header='{"alg":"RS256","typ":"JWT"}'
claim=$(printf '{
  "iss": "%s",
  "scope": "%s",
  "aud": "https://oauth2.googleapis.com/token",
  "exp": %s,
  "iat": %s
}' "$issuer" "$scopes" "$exp" "$iat")

encoded_header=$(echo "$header" | encode)
encoded_claim=$(echo "$claim" | encode)

printf '%s.%s' "${encoded_header}" "${encoded_claim}" >"$DIR/request"
printf '%s.%s.' "${encoded_header}" "${encoded_claim}" >"$DIR/token"

openssl pkcs12 -in "$p12" -out "$DIR/key" -nocerts -nodes -passin pass:notasecret
openssl dgst -sha256 -sign "$DIR/key" -out - "$DIR/request" | encode >>"$DIR/token"

assertion=$(cat "$DIR/token")

set -e
response=$(curl --silent --fail -d "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=$assertion" https://oauth2.googleapis.com/token)
echo "$response" | jq -r '.access_token'
