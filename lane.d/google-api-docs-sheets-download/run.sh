#!/bin/sh

if [ ! -x "$(command -v curl)" ]; then
  echo "Error: curl is not installed." >&2
  exit 3
fi

unset -v output id token
format='csv'
while getopts "i:t:f:o:" option; do
  case $option in
  o) output="$OPTARG" ;;
  i) id="$OPTARG" ;;
  t) token="$OPTARG" ;;
  f) format="$OPTARG" ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$id" ] || [ -z "$token" ] || [ -z "$output" ]; then
  echo "Must provide output path, document id and bearer token."
  exit 111
fi

echo 'Warning: this action is deprecated, see help for details' >&2
DIR=$(mktemp -dq)

trap 'set +x; rm -fr $DIR >/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

url=$(printf 'https://docs.google.com/spreadsheets/d/%s/export?exportFormat=%s' "$id" "$format")
header=$(printf 'Authorization: Bearer %s' "$token")

output_parent=$(dirname "$output")
mkdir -p "$output_parent"
curl --silent --fail -L -o "$output" --header "$header" "$url"
