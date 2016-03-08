package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

/*
Delete allows a user to make an outbound "DELETE /resource/:id" request.

	resp, err := jsh.Delete("http://apiserver", "user", "2")
*/
func Delete(urlStr string, resourceType string, id string) (*http.Response, error) {
	request, err := DeleteRequest(urlStr, resourceType, id)
	if err != nil {
		return nil, err
	}

	_, response, err := Do(request, jsh.ObjectMode)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/*
DeleteRequest returns a fully formatted request for performing a JSON API DELETE.
This is useful for if you need to set custom headers on the request. Otherwise
just use "jsc.Delete".
*/
func DeleteRequest(urlStr string, resourceType string, id string) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", err.Error()))
	}

	setIDPath(u, resourceType, id)

	request, err := NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error creating DELETE request: %s", err.Error()))
	}

	return request, nil
}
