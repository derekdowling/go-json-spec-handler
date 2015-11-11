package japi

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

const (
	// ContentType is the data encoding of choice for HTTP Request and Response Headers
	ContentType = "application/vnd.api+json"
)

// ParseObject returns a JSON object for a given io.ReadCloser containing
// a raw JSON payload
func ParseObject(reader io.ReadCloser) (*Object, error) {
	defer reader.Close()

	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request := struct {
		Object *Object `json:"data"`
	}{}

	err = json.Unmarshal(byteData, request)
	if err != nil {
		return nil, err
	}

	return request.Object, nil
}

// ParseList returns a JSON List for a given io.ReadCloser containing
// a raw JSON payload
func ParseList(reader io.ReadCloser) ([]*Object, error) {
	defer reader.Close()

	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request := struct {
		List []*Object `json:"data"`
	}{}

	err = json.Unmarshal(byteData, request)
	if err != nil {
		return nil, err
	}

	return request.List, nil
}
