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
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing URL: %s", err.Error())
	}

	setIDPath(u, object.Type, object.ID)

	request, err := http.NewRequest("PATCH", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating PATCH request: %s", err.Error())
	}

	return doObjectRequest(request, object)
}
