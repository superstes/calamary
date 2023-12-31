package parse

import (
	"fmt"
	"io"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
	tls_dissector "github.com/superstes/calamary/proc/parse/tls"
)

func parseTls(pkt ParsedPacket, conn net.Conn, connIo io.Reader, hdr [cnf.BYTES_HDR_L5]byte) (isTls meta.OptBool, tlsVersion uint16, sni string) {
	isTlsRaw := hdr[0] == tls_dissector.Handshake
	if isTlsRaw {
		isTls = meta.OptBoolTrue
	} else {
		isTls = meta.OptBoolFalse
	}

	if isTlsRaw {
		record, err := tls_dissector.ReadRecord(connIo)
		if err != nil {
			LogConnWarn("parse", pkt, fmt.Sprintf("Failed to parse TLS handshake: %v", err))
			return
		}

		clientHello := tls_dissector.ClientHelloMsg{}
		if err = clientHello.Decode(record.Opaque); err != nil {
			LogConnWarn("parse", pkt, fmt.Sprintf("Failed to parse TLS client-hello: %v", err))
			return
		}

		// todo: update connection TLS-version after serverHello
		if clientHello.Version != 0 {
			tlsVersion = uint16(clientHello.Version)
		} else {
			tlsVersion = uint16(record.Version)
		}

		for _, ext := range clientHello.Extensions {
			if ext.Type() == tls_dissector.ExtServerName {
				snExtension := ext.(*tls_dissector.ServerNameExtension)
				sni = snExtension.Name
				break
			}
		}
	}
	LogConnDebug("parse", pkt, fmt.Sprintf("TLS information: IsTls=%v, TlsVersion=%v, TlsSni=%s", isTlsRaw, tlsVersion, sni))
	return
}
