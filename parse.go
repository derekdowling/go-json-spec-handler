package japi

import (
	"encoding/json"
	"fmt"
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

	data := struct {
		Object Object `json:"data"`
	}{}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		)
	}

	return &data.Object, nil
}

// ParseList returns a JSON List for a given io.ReadCloser containing
// a raw JSON payload
func ParseList(reader io.ReadCloser) ([]*Object, error) {
	defer reader.Close()

	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	data := struct {
		List []*Object `json:"data"`
	}{List: []*Object{}}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		)
	}

	return data.List, nil
}
