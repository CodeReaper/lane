#!/bin/sh

TMP=$(mktemp -dq)
trap 'set +x; rm -rf $TMP 2>/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

unset -v type
unset -v input
unset -v key_row
unset -v main_language
unset -v output
while getopts "t:i:c:k:m:o:" option; do
  case $option in
  t) type=$OPTARG ;;
  i) input=$OPTARG ;;
  c) echo "$OPTARG" >>"$TMP/mapping" ;;
  k) key_row=$OPTARG ;;
  m) main_language=$OPTARG ;;
  o) output=$OPTARG ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ ! "$type" = "ios" ] && [ ! "$type" = "android" ]; then
  printf 'Provide ios or android as type.\n\n'
  exit 111
fi

if [ "$type" = "ios" ]; then
  if [ -z "$main_language" ]; then
    printf 'Provide main language.\n\n'
    exit 111
  fi
  if [ -z "$output" ]; then
    printf 'Provide output file.\n\n'
    exit 111
  fi
fi

if [ ! -f "$input" ]; then
  printf 'Provide input file.\n\n'
  exit 111
fi

if [ -z "$key_row" ]; then
  printf 'Provide key row.\n\n'
  exit 111
fi

makedir() {
  parent=$(dirname "$1")
  mkdir -p "$parent" 2>/dev/null
}

python=$(command -v python3 | head -n1)
if [ ! -x "$python" ]; then
  echo "Error: python3 is not installed." >&2
  exit 1
fi

cat >"${TMP}/extract" <<EOF
import sys
import csv
from xml.sax.saxutils import escape

key = int(sys.argv[1]) - 1
value = int(sys.argv[2]) - 1
should_escape = int(sys.argv[3]) == 1

with open(sys.stdin.fileno()) as file:
    reader = csv.reader(file, delimiter=',')
    for row in reader:
        if row[key]:
            v = row[value].replace("\n", "\\\\n")
            print("|".join([row[key], escape(v) if should_escape else v]))
EOF

while read -r item; do
  offset=$(echo "$item" | cut -d\  -f1 | tr -d "[:blank:]")
  [ "$type" = "android" ] && escape=1 || escape=0
  tail +2 "$input" | sed 's/|/\\\\\\/g' | $python "${TMP}/extract" "$key_row" "$offset" "$escape" | sed 's/\\\\\\/|/g' | LC_ALL=C sort -t \| -k1,1 >"${TMP}/${offset}.csv"
done <"${TMP}/mapping"

if [ "$type" = "ios" ]; then
  while read -r item; do
    offset=$(echo "$item" | cut -d\  -f1 | tr -d "[:blank:]")
    file=$(echo "$item" | cut -d\  -f2- | tr -d "[:blank:]")

    makedir "$file"

    printf "" >"$file"
    while read -r line; do
      key=$(printf "%s" "$line" | cut -d\| -f1 | sed 's|^"||g;s|"$||g')
      value=$(printf "%s" "$line" | cut -d\| -f2- | sed 's|^"||;s|"$||;s|\"|\\"|g')
      printf "\"%s\" = \"%s\";\n" "$key" "$value" >>"$file"
    done <"${TMP}/${offset}.csv"

  done <"${TMP}/mapping"

  makedir "$output"
  {
    echo '// swiftlint:disable all'
    echo 'import Foundation'
    echo 'struct Translations {'

    while read -r item; do
      key=$(printf "%s" "$item" | cut -d\| -f1 | sed 's|^"||g;s|"$||g')
      value=$(printf "%s" "$item" | cut -d\| -f2-)
      parameters=$(echo "$value" | grep -o -E '%[0-9]+' | wc -l | tr -d ' \n')

      if [ "$parameters" = "0" ]; then
        printf "\tstatic let %s = NSLocalizedString(\"%s\", comment: \"\")\n" "$key" "$key"
      else
        arguments=$(for i in $(seq 1 "$parameters"); do printf "p%s: String, _ " "$i"; done | rev | cut -c5- | rev)
        replacements=$(for i in $(seq 1 "$parameters"); do printf ".replacingOccurrences(of: \"%%%s\", with: p%s)" "$i" "$i"; done)
        printf "\tstatic func %s(_ %s) -> String {" "$key" "$arguments"
        printf " return NSLocalizedString(\"%s\", comment: \"\")" "$key"
        printf "%s" "$replacements"
        echo " }"
      fi
    done <"${TMP}/${main_language}.csv"

    echo '}'
  } >"$output"
fi

if [ "$type" = "android" ]; then
  while read -r item; do
    offset=$(echo "$item" | cut -d\  -f1 | tr -d "[:blank:]")
    file=$(echo "$item" | cut -d\  -f2- | tr -d "[:blank:]")

    makedir "$file"

    echo "<resources>" >"$file"
    while read -r line; do
      key=$(printf "%s" "$line" | cut -d\| -f1 | sed 's|^"||g;s|"$||g' | tr "[:upper:]" "[:lower:]")
      value=$(printf "%s" "$line" | cut -d\| -f2- | sed -E 's|%|%%|g;s|%%([0-9]+)|%\1$s|g;s|^"||;s|"$||')
      printf "\t<string name=\"%s\">%s</string>\n" "$key" "$value" >>"$file"
    done <"${TMP}/${offset}.csv"
    echo "</resources>" >>"$file"

  done <"${TMP}/mapping"
fi
