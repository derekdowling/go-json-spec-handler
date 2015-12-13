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
func Patch(urlStr string, object *jsh.Object) (*jsh.Object, *http.Response, *jsh.Error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", err.Error()))
	}

	setIDPath(u, object.Type, object.ID)

	request, err := http.NewRequest("PATCH", u.String(), nil)
	if err != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error creating PATCH request: %s", err.Error()))
	}

	return sendObjectRequest(request, object)
}
