package send

import (
	"io"
	"net"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func Forward(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	if pkt.L4.Proto == meta.ProtoL4Udp {
		// connFwd, err = DialUdp(pkt)
		parse.LogConnError("send", pkt, "UDP not yet implemented!")

	} else {
		forwardTcp(pkt, conn, connIo)
	}
}
