package jsh

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// CreateReadCloser is a helper function for dealing with creating HTTP requests
func CreateReadCloser(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	return ioutil.NopCloser(reader)
}

func testRequest(bytes []byte) (*http.Request, error) {
	req, err := http.NewRequest("GET", "", CreateReadCloser(bytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ContentType)
	return req, nil
}
