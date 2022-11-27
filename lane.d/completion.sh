_lane() {
  COMPREPLY=()
  if [ "$1" = "$3" ]; then
    if [ "$2" = "--" ]; then
      COMPREPLY=($(lane lanes | tr "\n" ' '))
    else
      COMPREPLY=($(lane lanes | grep "^$2" | tr "\n" ' '))
    fi
  fi
  return 0
}

complete -F _lane lane
