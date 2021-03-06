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

# awk too magical
sed 's/|/\\\\\\/g' <"$input" | awk -vOFS='|' -vq='"' 'func csv2del(n) { for(i=n; i<=c; i++) {if(i%2 == 1) gsub(/,/, OFS, a[i]); else a[i] = (q a[i] q); out = (out) ? out a[i] : a[i]}; return out} {c=split($0, a, q); out=X; if(a[1]) $0=csv2del(1); else $0=csv2del(2)}1' >"${TMP}/input"

while read -r item; do
  offset=$(echo "$item" | cut -d\  -f1 | tr -d "[:blank:]")
  tail +2 "${TMP}/input" | grep -v ^$ | cut -d\| -f"$key_row,$offset" | sed 's/\\\\\\/|/g' | sort >"${TMP}/${offset}.csv"
done <"${TMP}/mapping"

if [ "$type" = "ios" ]; then
  while read -r item; do
    offset=$(echo "$item" | cut -d\  -f1 | tr -d "[:blank:]")
    file=$(echo "$item" | cut -d\  -f2- | tr -d "[:blank:]")

    makedir "$file"

    printf "" >"$file"
    while read -r line; do
      key=$(echo "$line" | cut -d\| -f1 | sed 's|^"||g;s|"$||g')
      value=$(echo "$line" | cut -d\| -f2- | sed 's|^"||;s|"$||;s|\"|\\"|g')
      echo "\"$key\" = \"$value\";" >>"$file"
    done <"${TMP}/${offset}.csv"

  done <"${TMP}/mapping"

  makedir "$output"
  {
    echo '// swiftlint:disable all'
    echo 'import Foundation'
    echo 'struct Translations {'

    while read -r item; do
      key=$(echo "$item" | cut -d\| -f1 | sed 's|^"||g;s|"$||g')
      value=$(echo "$item" | cut -d\| -f2-)
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
      key=$(echo "$line" | cut -d\| -f1 | sed 's|^"||g;s|"$||g' | tr "[:upper:]" "[:lower:]")
      value=$(echo "$line" | cut -d\| -f2- | sed -E 's|(%[0-9]+)|\1$s|g;s|^"||;s|"$||')
      printf "\t<string name=\"%s\">%s</string>\n" "$key" "$value" >>"$file"
    done <"${TMP}/${offset}.csv"
    echo "</resources>" >>"$file"

  done <"${TMP}/mapping"
fi
