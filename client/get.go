package jsc

import (
	"fmt"
	"net/http"
	"net/url"
)

// Get allows a user to make an outbound GET /resources(/:id) request.
//
// For a GET request that retrieves multiple resources, pass an empty string for
// the id parameter:
//
//  GET "http://apiserver/users
//	resp, err := jsh.Get("http://apiserver", "user", "")
//	list, err := resp.GetList()
//
// For a GET request on a specific attribute:
//
//  GET "http://apiserver/users/2
//	resp, err := jsh.Get("http://apiserver", "user", "2")
//	obj := resp.GetObject()
//
func Get(urlStr string, resourceType string, id string) (*Response, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// ghetto pluralization, fix when it becomes an issue
	setPath(u, resourceType)

	response, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("Error performing GET request: %s", err.Error())
	}

	return &Response{response}, nil
}
