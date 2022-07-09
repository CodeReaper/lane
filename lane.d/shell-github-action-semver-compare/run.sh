#!/bin/sh

unset -v MAIN
unset -v CURRENT
while getopts "m:c:" option; do
  case $option in
  m) MAIN=$OPTARG ;;
  c) CURRENT=$OPTARG ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$MAIN" ] || [ -z "$CURRENT" ]; then
  printf 'Provide both main and current versions.\n\n'
  exit 111
fi

echo '::group::Resolved versions'
printf '%s - main version\n%s - current version\n' "$MAIN" "$CURRENT"
echo '::endgroup::'

[ "$MAIN" = "$CURRENT" ] && {
  echo 'Version must be changed.'
  exit 10
}

VERIFY=$(printf '%s\n%s' "$MAIN" "$CURRENT" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -k 4,4 -k 5,5 -g | tail -n1)
[ "$MAIN" = "$VERIFY" ] && {
  echo 'Version must be greater than version on main.'
  exit 20
}

exit 0
