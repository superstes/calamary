#!/bin/bash

PROXY_HOST='172.17.1.81'
PROXY_USER='tester'
PROXY_SSH_PORT=22

function ssh_cmd {
  ssh -p "$PROXY_SSH_PORT" "$PROXY_USER"@"$PROXY_HOST" "$1"
}