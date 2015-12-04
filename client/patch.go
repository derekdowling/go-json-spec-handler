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
//  obj, _ := jsh.NewObject("123", "resource_name", payload)
//  resp, _ := jsc.Patch("http://postap.com", obj)
//  updatedObj, _ := resp.GetObject()
//
func Patch(urlStr string, object *jsh.Object) (*Response, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	setIDPath(u, object.Type, object.ID)

	request, err := http.NewRequest("PATCH", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating new HTTP request: %s", err.Error())
	}

	return sendObjectRequest(request, object)
}
