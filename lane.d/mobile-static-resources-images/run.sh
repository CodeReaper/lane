#!/bin/sh

unset -v input
unset -v output
while getopts "i:o:" option; do
  case $option in
  i) input=$OPTARG ;;
  o) output=$OPTARG ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$input" ] || [ -z "$output" ]; then
  printf 'Provide both input and output arguments.\n\n'
  exit 111
fi

if [ ! -d "$input" ]; then
  printf 'Provide a valid input directory.\n\n'
  exit 3
fi

parent=$(dirname "$output")
mkdir -p "$parent" 2>/dev/null

{
  echo '// swiftlint:disable all'
  echo 'import UIKit'
  echo 'struct Images {'

  find "$input" -type d -iname "*.imageset" | while read -r item; do
    name=$(basename "$item" .imageset)
    safe_name=$(echo "$name" | sed 's/-/_/g;s/ /_/g')
    printf "\tstatic let %s = UIImage(named:\"%s\")!\n" "$safe_name" "$name"
  done

  echo '}'
} >"$output"
