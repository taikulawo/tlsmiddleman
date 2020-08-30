package common

import (
	"bufio"
	"bytes"
	"io"
)

func NewSliceWriter() *bytes.Buffer {
	return &bytes.Buffer{}
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(r interface{}, err error) interface{} {
	Must(err)
	return r
}

func ReaderAndWriter(rw io.ReadWriter) (c1,c2 chan[]byte) {
	c1 := make(chan []byte,1)
	c2 := make(chan []byte,1)
	bufio.NewWriter()
	bufio.NewReader()
	go func() {
		rw.Write(<- c1)
	}()
	go func() {
		_, err := rw.Read(<- c2)
		if err != nil {
		}
	}()
	return c1, c2
}
