#!/bin/bash

set -uo pipefail
set +e

PORT_BASE="${PORT_BASE:-1000}"
PROXY_PORT="${PROXY_PORT:-"${PORT_BASE}1"}"
# tests are targeting the 'transparent' mode

export PROXY_HOST="$PROXY_HOST"
export PROXY_PORT="$PROXY_PORT"

testsGeneral=()

# log_header 'RUNNING TESTS: GENERAL'

for test in "${testsGeneral[@]}"
do
  if ! runTest "test_general/$test"
  then
    fail
  fi
done

unset PROXY_PORT
