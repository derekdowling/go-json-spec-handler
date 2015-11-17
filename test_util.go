package jsh

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func createIOCloser(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	return ioutil.NopCloser(reader)
}

func testRequest(bytes []byte) (*http.Request, error) {
	req, err := http.NewRequest("GET", "", createIOCloser(bytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ContentType)
	return req, nil
}
