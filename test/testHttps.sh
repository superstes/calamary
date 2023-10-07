#!/bin/bash

set -uo pipefail
set +e

source ./util/base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}5"}"

export http_proxy="https://${PROXY_HOST}:${PROXY_PORT}"
export https_proxy="https://${PROXY_HOST}:${PROXY_PORT}"
export HTTP_PROXY="https://${PROXY_HOST}:${PROXY_PORT}"
export HTTPS_PROXY="https://${PROXY_HOST}:${PROXY_PORT}"

testsHttps=()

# log_header 'RUNNING TESTS: HTTPS'

for test in "${testsHttps[@]}"
do
  if ! runTest "https/$test"
  then
    fail
  fi
done

unset PROXY_PORT
unset http_proxy
unset https_proxy
unset HTTP_PROXY
unset HTTPS_PROXY
