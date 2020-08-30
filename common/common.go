package common

import (
	"bytes"
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
