package u

import (
	"bufio"
	"io"
	"net"
)

const (
	bufferSize = 64 * 1024
)

func Transport(rw1, rw2 io.ReadWriter) error {
	errc := make(chan error, 1)
	go func() {
		errc <- CopyBuffer(rw1, rw2, bufferSize)
	}()

	go func() {
		errc <- CopyBuffer(rw2, rw1, bufferSize)
	}()

	if err := <-errc; err != nil && err != io.EOF {
		return err
	}

	return nil
}

func CopyBuffer(dst io.Writer, src io.Reader, bufSize int) error {
	buf := GetBufferPool(bufSize)
	defer PutBufferPool(buf)

	_, err := io.CopyBuffer(dst, src, *buf)
	return err
}

type bufferReaderConn struct {
	net.Conn
	br *bufio.Reader
}
