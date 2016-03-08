package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

// Fetch performs an outbound GET /resourceTypes/:id request
func Fetch(baseURL string, resourceType string, id string) (*jsh.Document, *http.Response, error) {
	request, err := FetchRequest(baseURL, resourceType, id)
	if err != nil {
		return nil, nil, err
	}

	return Do(request, jsh.ObjectMode)
}

/*
FetchRequest returns a fully formatted JSONAPI Fetch request. Useful if you need to
set custom headers before proceeding. Otherwise just use "jsh.Fetch".
*/
func FetchRequest(baseURL string, resourceType, id string) (*http.Request, error) {
	if id == "" {
		return nil, jsh.SpecificationError("ID cannot be empty for GetObject request type")
	}

	u, urlErr := url.Parse(baseURL)
	if urlErr != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setIDPath(u, resourceType, id)

	return NewRequest("GET", u.String(), nil)
}

// List prepares an outbound GET /resourceTypes request
func List(baseURL string, resourceType string) (*jsh.Document, *http.Response, error) {
	request, err := ListRequest(baseURL, resourceType)
	if err != nil {
		return nil, nil, err
	}

	return Do(request, jsh.ListMode)
}

/*
ListRequest returns a fully formatted JSONAPI List request. Useful if you need to
set custom headers before proceeding. Otherwise just use "jsh.List".
*/
func ListRequest(baseURL string, resourceType string) (*http.Request, error) {
	u, urlErr := url.Parse(baseURL)
	if urlErr != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setPath(u, resourceType)

	return NewRequest("GET", u.String(), nil)
}
