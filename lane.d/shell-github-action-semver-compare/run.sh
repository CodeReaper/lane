#!/bin/sh

unset -v main
unset -v current
quiet=0
while getopts "m:c:q" option; do
  case $option in
  m) main=$OPTARG ;;
  c) current=$OPTARG ;;
  q) quiet=1 ;;
  \?) exit 111 ;;
  esac
done
shift $((OPTIND - 1))

if [ -z "$main" ] || [ -z "$current" ]; then
  printf 'Provide both main and current versions.\n\n'
  exit 111
fi

if [ $quiet -eq 0 ]; then
  echo '::group::Resolved versions'
  printf '%s - main version\n%s - current version\n' "$main" "$current"
  echo '::endgroup::'
fi

[ "$main" = "$current" ] && {
  echo 'Version must be changed.'
  exit 10
}

verify=$(printf '%s\n%s' "$main" "$current" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -k 4,4 -k 5,5 -g | tail -n1)
[ "$main" = "$verify" ] && {
  echo 'Version must be greater than version on main.'
  exit 20
}

exit 0
