package replicant

import (
	"bufio"
	"io"
	"net"
)

type Replicant struct {
	conn net.Conn
	remote net.Conn
}

func (this *Replicant) DumpRequest(from, to net.Conn, c <- chan []byte) *bufio.Reader {
	reader, writer := io.Pipe()
	go func() {
		for {
			b := make([]byte, 1024)
			from.Read(b)
			writer.Write(b)
		}
	}()
	return bufio.NewReader(reader)
}