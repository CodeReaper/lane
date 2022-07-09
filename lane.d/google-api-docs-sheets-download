#!/bin/sh

if [ ! -x "$(command -v curl)" ]; then
  echo "Error: curl is not installed." >&2
  exit 3
fi

usage() {
cat << END
OPTIONS
    -h show this usage
    -o output path
    -i document id
    -t bearer token, see google-api-jwt-generate
    -f format, defaults to csv
END
}

unset -v output id token
format='csv'
while getopts "hi:t:f:o:" option; do
  case $option in
    o) output="$OPTARG" ;;
    i) id="$OPTARG" ;;
    t) token="$OPTARG" ;;
    f) format="$OPTARG" ;;
    \?)
      echo "unknown option: $option"
      usage
      exit 1
      ;;
    h)
      usage
      exit 0
      ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$id" ] || [ -z "$token" ] || [ -z "$output" ]; then
    echo "Must provide output path, document id and bearer token."
    usage
    exit 2
fi

DIR=$(mktemp -dq)

trap 'set +x; rm -fr $DIR >/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

url=$(printf 'https://docs.google.com/spreadsheets/d/%s/export?exportformat=%s' "$id" "$format")
header=$(printf 'Authorization: Bearer %s' "$token")

output_parent=$(dirname "$output")
mkdir -p "$output_parent"
curl --silent --fail -L -o "$output" --header "$header" "$url"
