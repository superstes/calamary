#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")"

PROXY_HOST='172.17.1.56'
PROXY_USER='proxy_test'
PROXY_SSH_PORT=22
TMP_BASE="/tmp/calamary_$(date +%s)"
PORT_BASE="$(date +'%H%M')"

# remove leading 0 as it is not valid as port
if [[ ${PORT_BASE:0:1} == "0" ]]
then
  PORT_BASE="3${PORT_BASE:1}"
fi

function log {
  echo ''
  echo "$1"
  echo ''
}

log 'PREPARING BINARY'
git clone

log 'PREPARING FOR TESTS'

cp 'config.yml' 'config_tmp.yml'

sed -i "s@PORT_BASE@$PORT_BASE@g" 'config_tmp.yml'
sed -i "s@CRT_BASE@$TMP_BASE@g" 'config_tmp.yml'

log 'GENERATING CERTS'
# todo: generate ca & subca
openssl req -x509 -newkey rsa:4096 -keyout 'cert_tmp.key' -out 'cert_tmp.crt' -sha256 -days 60 -nodes -subj "$CERT_CN"

log 'COPYING FILES TO PROXY-HOST'
function copy_file {
    scp -P "$PROXY_SSH_PORT" "$1" "$PROXY_USER"@"$PROXY_HOST":"$2"
}

copy_file 'config_tmp.yml' "${TMP_BASE}.yml'"
copy_file 'cert_tmp.key' "${TMP_BASE}.key"
copy_file 'cert_tmp.crt' "${TMP_BASE}.crt"

function runTest {
  testScript="$1"
  echo ''
  echo "RUNNING TEST '${testScript}'"
  echo ''
  ./${testScript}.sh
  result="$?"
  if [[ "result" -ne "0" ]]
  then
    echo "FAILED TEST '${testScript}'"
    return 1
  fi
  echo ''
  return 0
}

function fail {
  log 'TEST-RUN FAILED!'
  exit 1
}

testsToRun=()
testsToRun[0]="dummyOk"
testsToRun[1]="dummyFail"


log 'STARTING TESTS'

sed +e
source testTransparent.sh

log 'TEST-RUN FINISHED SUCCESSFULLY!'

exit 0
