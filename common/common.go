package common

import (
	"io"
	"net"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(r interface{}, err error) interface{} {
	Must(err)
	return r
}

// ReadHelper
type ReaderHelper struct {
	input <-chan []byte
}

func (this *ReaderHelper) Read(b []byte) (int, error) {
	m := <-this.input
	if m == nil {
		return 0, io.EOF
	}
	copy(b, m)
	return len(m), nil
}

func NewReaderHelper(c <-chan []byte) *ReaderHelper {
	return &ReaderHelper{
		input: c,
	}
}

// 从Conn获取channel
func ChannelFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				sent := make([]byte, n)
				copy(sent, b[:n])
				c <- sent
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()
	return c
}
