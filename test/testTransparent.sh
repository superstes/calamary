#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}1"}"

testsTransparent=()
testsTransparent[0]="basic"
testsTransparent[1]="dummyOk"
#testsTransparent[2]="dummyFail"

if [ "${#testsTransparent}" -gt "0" ]
then
  log_header 'RUNNING TESTS: TRANSPARENT'
fi

for test in "${testsTransparent[@]}"
do
  if ! runTest "transparent/$test"
  then
    fail
  fi
done

unset PROXY_PORT
