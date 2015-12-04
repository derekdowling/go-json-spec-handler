package jsc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

// Post allows a user to make an outbound POST /resources request:
//
//	obj, _ := jsh.NewObject("123", "user", payload)
//	resp, _ := jsh.Post("http://apiserver", obj)
//	createdObj := resp.GetObject()
//
func Post(urlStr string, object *jsh.Object) (*Response, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// ghetto pluralization, fix when it becomes an issue
	setPath(u, object.Type)

	request, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building POST request: %s", err.Error())
	}

	return sendObjectRequest(request, object)
}
