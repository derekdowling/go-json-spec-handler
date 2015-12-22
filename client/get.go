package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

// Fetch performs an outbound GET /resourceTypes/:id request
func Fetch(urlStr string, resourceType string, id string) (*jsh.Document, *http.Response, *jsh.Error) {
	if id == "" {
		return nil, nil, jsh.SpecificationError("ID cannot be empty for GetObject request type")
	}

	u, urlErr := url.Parse(urlStr)
	if urlErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setIDPath(u, resourceType, id)

	return Get(u.String())
}

// List prepares an outbound GET /resourceTypes request
func List(urlStr string, resourceType string) (*jsh.Document, *http.Response, *jsh.Error) {
	u, urlErr := url.Parse(urlStr)
	if urlErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setPath(u, resourceType)

	return Get(u.String())
}

// Get performs a generic GET request for a given URL and attempts to parse the
// response into a JSON API Format
func Get(urlStr string) (*jsh.Document, *http.Response, *jsh.Error) {
	response, httpErr := http.Get(urlStr)
	if httpErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error performing GET request: %s", httpErr.Error()))
	}

	doc, err := Document(response)
	if err != nil {
		return nil, nil, err
	}

	return doc, response, nil
}
