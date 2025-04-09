package conn

import (
	"bufio"
	"net"
)

type WrappedConn struct {
	net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func Wrap(c net.Conn) *WrappedConn {
	return &WrappedConn{
		Conn:   c,
		Reader: bufio.NewReader(c),
		Writer: bufio.NewWriter(c),
	}
}

func (wc *WrappedConn) ReadMessage() ([]byte, error) {
	return wc.Reader.ReadBytes('\n')
}

func (wc *WrappedConn) WriteMessage(data []byte) error {
	if _, err := wc.Writer.Write(append(data, '\n')); err != nil {
		return err
	}
	return wc.Writer.Flush()
}
