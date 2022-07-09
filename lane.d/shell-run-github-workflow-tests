#!/bin/sh

for command in git yq jq; do
  if ! [ -x "$(command -v $command)" ]; then
    echo "Error: $command is not installed." >&2
    exit 1
  fi
done

_USAGE() {
cat << END
OPTIONS
    -h show this usage
    -i workflow file
    -j job id
END
}

unset -v FILE
unset -v JOBID
while getopts "hi:j:" option; do
    case $option in
        i) FILE="$OPTARG" ;;
        j) JOBID="$OPTARG" ;;
        \?)
          echo "unknown option: $option"
          _USAGE
          exit 1
          ;;
        h)
          _USAGE
          exit 0
          ;;
    esac
done
shift $((OPTIND - 1))

if [ ! -f "$FILE" ]; then
    echo "Must provide a valid path to a GitHub workflow file with script-based tests."
    if [ -n "$FILE" ]; then
      printf "Given: %s\n\n" "$FILE"
    else
      echo
    fi
    _USAGE
    exit 1
fi

DIR=$(mktemp -dq)

trap 'set +x; rm -fr $DIR >/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

set -e

printf 'Preparing runnner...'

echo 'TOTAL_PASS=0; TOTAL_FAIL=0' > "$DIR/runner.sh"

yq -o json "$FILE" | jq -rc '.jobs | to_entries[] | [{group: .key, groupName: .value.name, script: .value.steps}] | .[]' | while read -r JOB; do
    GROUP=$(printf '%s\n' "$JOB" | jq -rj '.group')
    GROUP_NAME=$(printf '%s\n' "$JOB" | jq -rj '.groupName')

    {
      echo "cd '${DIR}/workspace/'; git reset HEAD --hard > /dev/null; git clean -fdx . > /dev/null; sh '${DIR}/${GROUP}.sh'"
      echo "TOTAL_PASS=\$((TOTAL_PASS+\$(cat ${DIR}/PASS))); TOTAL_FAIL=\$((TOTAL_FAIL+\$(cat ${DIR}/FAIL)))"
    } >> "$DIR/runner.sh"

    {
      echo "echo; echo '$GROUP_NAME ($GROUP)'"
      echo "GREEN=$'\e[0;32m'; RED=$'\e[0;31m'; NC=$'\e[0m'; PASS=0; FAIL=0"
    } > "${DIR}/${GROUP}.sh"

    if [ -n "$JOBID" ] && [ ! "$JOBID" = "$GROUP" ]; then
      {
        echo "echo ' - Skipped'"
        echo "printf 0 > '${DIR}/PASS'; printf 0 > '${DIR}/FAIL'"
      } >> "${DIR}/${GROUP}.sh"
      continue
    fi

    I=0
    printf '%s\n' "$JOB" | jq -rc '.script[] | select(.run != null)' | while read -r STEP; do
        STEP_NAME=$(printf '%s\n' "$STEP" | jq -rj '.name')
        RUN=$(printf '%s\n' "$STEP" | jq -rj '.run')

        echo "$RUN" > "${DIR}/${GROUP}.${I}.sh"

        {
          echo "printf ' - ${STEP_NAME}: '"
          echo "set +e; sh '${DIR}/${GROUP}.${I}.sh' > messages 2>&1"
          printf "if [ \$? -eq 0 ]; then printf \${GREEN}'Pass\n'\${NC}; PASS=\$((PASS+1)); else printf \${RED}'Fail\n'\${NC}; FAIL=\$((FAIL+1)); cat messages; fi;\n"
        } >> "${DIR}/${GROUP}.sh"

        I=$((I+1))
    done

    echo "printf \$PASS > '${DIR}/PASS'; printf \$FAIL > '${DIR}/FAIL'" >> "${DIR}/${GROUP}.sh"

done

{
  # shellcheck disable=SC2016
  echo 'TOTAL=$((TOTAL_PASS+TOTAL_FAIL))'
  # shellcheck disable=SC2016
  printf 'echo; printf "Tests; Total: \033[1m${TOTAL}\033[0m Passes: \033[1m${TOTAL_PASS}\033[0m Fails: \033[1m${TOTAL_FAIL}\033[0m\n"\n'
} >> "$DIR/runner.sh"

echo ' done!'

printf 'Preparing workspace... '
mkdir "$DIR/workspace"
tar -c --exclude .git . | tar -x -C "$DIR/workspace/"

cd "$DIR/workspace/"
git init > /dev/null 2>&1
git add . > /dev/null
git commit -am "known state" > /dev/null
echo ' done!'

trap 'set +x; cd - > /dev/null; rm -fr $DIR >/dev/null 2>&1' 0
trap 'exit 2' 1 2 3 15

echo 'Executing runner...'
sh "$DIR/runner.sh"
