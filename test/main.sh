#!/bin/bash

set -euo pipefail

source ./target.sh

TMP_BASE="/tmp/calamary_${VERSION}"  # could be problematic; but so will nat
PORT_BASE="$(date +'%H%M')"
CERT_CN="/C=AT/ST=Styria/CN=Calamary Forward Proxy"

# remove leading 0 as it is not valid as port
if [[ ${PORT_BASE:0:1} == "0" ]]
then
  PORT_BASE="3${PORT_BASE:1}"
fi

source ./util/main.sh

cleanup

log 'PREPARING FOR TESTS'

cp 'config.yml' 'config_tmp.yml'

sed -i "s@PORT_BASE@$PORT_BASE@g" 'config_tmp.yml'
sed -i "s@CRT_BASE@$TMP_BASE@g" 'config_tmp.yml'

log 'GENERATING CERTS'
easyrsa="$(pwd)/easyrsa/easyrsa"
export EASYRSA_PKI="$(pwd)/pki"
mkdir -p "$EASYRSA_PKI"
export EASYRSA_REQ_COUNTRY='AT'
export EASYRSA_REQ_PROVINCE='Styria'
export EASYRSA_CERT_EXPIRE='60'
export PKI_CERT_CN='Test CA'

easyrsa init-pki
easyrsa build-ca nopass
export PKI_CERT_CN='Test Server'
easyrsa gen-req server nopass
easyrsa sign-req server proxy

export EASYRSA_PKI="$(pwd)/pki_sub"
mkdir -p "$EASYRSA_PKI"
export PKI_CERT_CN='Test SSL-Interception'
easyrsa init-pki
easyrsa build-ca nopass subca

export EASYRSA_PKI="$(pwd)/pki"
easyrsa import-req "$(pwd)/pki_sub/reqs/ca.req" subca
easyrsa sign-req ca subca

chmod 644 -R "$(pwd)/"pki*

log 'COPYING FILES TO PROXY-HOST'

copy_file 'calamary' "$TMP_BASE"
copy_file 'config_tmp.yml' "${TMP_BASE}.yml"
copy_file "$(pwd)/pki/ca.crt" "${TMP_BASE}.ca.crt"
copy_file "$(pwd)/pki/issued/subca.crt" "${TMP_BASE}.subca.crt"
copy_file "$(pwd)/pki_sub/private/ca.key" "${TMP_BASE}.subca.key"
copy_file "$(pwd)/pki/issued/proxy.crt" "${TMP_BASE}.crt"
copy_file "$(pwd)/pki/private/proxy.key" "${TMP_BASE}.key"

ssh_cmd "sudo chown proxy:proxy ${TMP_BASE}*"

log 'STARTING PROXY'
ssh_cmd "sudo systemctl start calamary@${VERSION}.service"

# todo: iptables/nftables NAT for transparent mode

log 'STARTING TESTS'

sed +e
find . -type f -name '*.sh' -exec chmod +x {} \;
source testTransparent.sh
source testGeneral.sh
source testTransparentTproxy.sh
source testHttp.sh
source testHttps.sh
source testProxyproto.sh
source testSocks5.sh

log 'TEST-RUN FINISHED SUCCESSFULLY!'
status='PASSED'

cleanup
update_badge

exit 0
