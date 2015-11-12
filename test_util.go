package jsh

import (
	"bytes"
	"io"
	"io/ioutil"
)

func createIOCloser(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	return ioutil.NopCloser(reader)
}
