package jsh

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

const (
	// ContentType is the data encoding of choice for HTTP Request and Response Headers
	ContentType = "application/vnd.api+json"
)

// ParseObject returns a JSON object for a given io.ReadCloser containing
// a raw JSON payload
func ParseObject(reader io.ReadCloser) (*Object, SendableError) {
	defer closeReader(reader)

	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Error attempting to read request body: %s", err))
	}

	data := struct {
		Object Object `json:"data"`
	}{}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		))
	}

	object := &data.Object
	return object, validateInput(object)
}

// ParseList returns a JSON List for a given io.ReadCloser containing
// a raw JSON payload
func ParseList(reader io.ReadCloser) ([]*Object, SendableError) {
	defer closeReader(reader)

	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Error attempting to read request body: %s", err))
	}

	data := struct {
		List []*Object `json:"data"`
	}{List: []*Object{}}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		))
	}

	for _, object := range data.List {
		err := validateInput(object)
		if err != nil {
			return nil, err
		}
	}

	return data.List, nil
}

func closeReader(reader io.ReadCloser) {
	err := reader.Close()
	if err != nil {
		log.Println("Unabled to close request Body")
	}
}
