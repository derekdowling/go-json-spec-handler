package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

// GetObject allows a user to make an outbound GET /resourceTypes/:id
func GetObject(urlStr string, resourceType string, id string) (*jsh.Object, *http.Response, *jsh.Error) {
	if id == "" {
		return nil, nil, jsh.SpecificationError("ID cannot be empty for GetObject request type")
	}

	u, urlErr := url.Parse(urlStr)
	if urlErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setIDPath(u, resourceType, id)

	response, err := Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	object, err := ParseObject(response)
	if err != nil {
		return nil, response, err
	}

	return object, response, nil
}

// GetList prepares an outbound request for /resourceTypes expecting a list return value.
func GetList(urlStr string, resourceType string) (jsh.List, *http.Response, *jsh.Error) {
	u, urlErr := url.Parse(urlStr)
	if urlErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setPath(u, resourceType)

	response, err := Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	list, err := ParseList(response)
	if err != nil {
		return nil, response, err
	}

	return list, response, nil
}

// Get performs a Get request for a given URL and returns a basic Response type
func Get(urlStr string) (*http.Response, *jsh.Error) {
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error performing GET request: %s", err.Error()))
	}

	return response, nil
}
