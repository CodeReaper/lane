#!/bin/sh

TMP=$(mktemp -dq)
trap 'set +x; rm -rf $TMP 2>/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

unset -v TYPE
unset -v INPUT
unset -v KEY_ROW
unset -v MAIN_LANGUAGE
unset -v OUTPUT
unset -v BETTER
while getopts "t:i:c:k:m:o:" option; do
  case $option in
  c) echo "$OPTARG" >>"$TMP/mapping" ;;
  t) TYPE=$OPTARG ;;
  i) INPUT=$OPTARG ;;
  k) KEY_ROW=$OPTARG ;;
  m) MAIN_LANGUAGE=$OPTARG ;;
  o) OUTPUT=$OPTARG ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ ! "$TYPE" = "ios" ] && [ ! "$TYPE" = "android" ]; then
  printf 'Provide ios or android as type.\n\n'
  exit 111
fi

if [ "$TYPE" = "ios" ]; then
  if [ -z "$MAIN_LANGUAGE" ]; then
    printf 'Provide main language.\n\n'
    exit 111
  fi
  if [ -z "$OUTPUT" ]; then
    printf 'Provide output file.\n\n'
    exit 111
  fi
fi

if [ ! -f "$INPUT" ]; then
  printf 'Provide input file.\n\n'
  exit 111
fi

if [ -z "$KEY_ROW" ]; then
  printf 'Provide key row.\n\n'
  exit 111
fi

_MKDIR() {
  DIRECTORY=$(dirname "$1")
  mkdir -p "$DIRECTORY" 2>/dev/null
}

_UNPACK_CSV() {
  while read -r ITEM; do
    OFFSET=$(echo "$ITEM" | cut -d\  -f1 | tr -d "[:blank:]")
    tail +2 "./$INPUT" | grep -v ^$ | sed 's|\\;|\\\\\\|g' | cut -d\; -f"$KEY_ROW,$OFFSET" | sed 's|\\\\\\|;|g' | sort >"${TMP}/${OFFSET}.csv"
  done <"${TMP}/mapping"
}

_GENERATE_XML() {
  while read -r ITEM; do
    OFFSET=$(echo "$ITEM" | cut -d\  -f1 | tr -d "[:blank:]")
    FILE=$(echo "$ITEM" | cut -d\  -f2- | tr -d "[:blank:]")

    _MKDIR "$FILE"

    echo "<resources>" >"$FILE"
    while read -r LINE; do
      KEY=$(echo "$LINE" | cut -d\; -f1 | tr "[:upper:]" "[:lower:]")
      VALUE=$(echo "$LINE" | cut -d\; -f2- | sed -E 's|(%[0-9]+)|\1$s|g')
      printf "\t<string name=\"%s\">%s</string>\n" "$KEY" "$VALUE" >>"$FILE"
    done <"${TMP}/${OFFSET}.csv"
    echo "</resources>" >>"$FILE"

  done <"${TMP}/mapping"
}

_GENERATE_STRINGS() {
  while read -r ITEM; do
    OFFSET=$(echo "$ITEM" | cut -d\  -f1 | tr -d "[:blank:]")
    FILE=$(echo "$ITEM" | cut -d\  -f2- | tr -d "[:blank:]")

    _MKDIR "$FILE"

    printf "" >"$FILE"
    while read -r LINE; do
      KEY=$(echo "$LINE" | cut -d\; -f1)
      VALUE=$(echo "$LINE" | cut -d\; -f2- | sed 's|\"|\\"|g')
      echo "\"$KEY\" = \"$VALUE\";" >>"$FILE"
    done <"${TMP}/${OFFSET}.csv"

  done <"${TMP}/mapping"
}

_GENERATE_STRUCT() {
  echo '// swiftlint:disable all'
  echo 'import Foundation'
  echo 'struct Translations {'

  while read -r ITEM; do
    KEY=$(echo "$ITEM" | cut -d\; -f1)
    VALUE=$(echo "$ITEM" | cut -d\; -f2-)
    PARAMETERS=$(echo "$VALUE" | grep -o -E '%[0-9]+' | wc -l | tr -d ' \n')

    if [ "$PARAMETERS" = "0" ]; then
      printf "\tstatic let %s = NSLocalizedString(\"%s\", comment: \"\")\n" "$KEY" "$KEY"
    else
      ARGUMENTS=$(for i in $(seq 1 "$PARAMETERS"); do printf "p%s: String, _ " "$i"; done | rev | cut -c5- | rev)
      REPLACEMENTS=$(for i in $(seq 1 "$PARAMETERS"); do printf ".replacingOccurrences(of: \"%%%s\", with: p%s)" "$i" "$i"; done)
      printf "\tstatic func %s(_ %s) -> String {" "$KEY" "$ARGUMENTS"
      printf " return NSLocalizedString(\"%s\", comment: \"\")" "$KEY"
      printf "%s" "$REPLACEMENTS"
      echo " }"
    fi
  done <"${TMP}/${MAIN_LANGUAGE}.csv"

  echo '}'
}

_UNPACK_CSV

if [ "$TYPE" = "ios" ]; then
  _GENERATE_STRINGS

  _MKDIR "$OUTPUT"
  _GENERATE_STRUCT >"$OUTPUT"
fi

if [ "$TYPE" = "android" ]; then
  _GENERATE_XML
fi
