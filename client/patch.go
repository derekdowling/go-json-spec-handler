package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

// Patch allows a consumer to perform a PATCH /resources/:id request
// Example:
//
//  obj, _ := jsh.NewObject("123", "user", payload)
//	// does PATCH /http://postap.com/api/user/123
//  json, resp, err := jsc.Patch("http://postap.com/api/", obj)
//	updatedObj := json.First()
//
func Patch(baseURL string, object *jsh.Object) (*jsh.Document, *http.Response, error) {
	request, err := PatchRequest(baseURL, object)
	if err != nil {
		return nil, nil, err
	}

	return Do(request)
}

// PatchRequest returns a fully formatted request with JSON body for performing
// a JSONAPI PATCH. This is useful for if you need to set custom headers on the
// request. Otherwise just use "jsc.Patch".
func PatchRequest(baseURL string, object *jsh.Object) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL: %s", err.Error())
	}

	setIDPath(u, object.Type, object.ID)

	request, err := NewRequest("PATCH", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating PATCH request: %s", err.Error())
	}

	err = prepareBody(request, object)
	if err != nil {
		return nil, err
	}

	return request, nil
}
