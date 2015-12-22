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
//	// does POST http://apiserver/user/123
//	json, resp, err := jsh.Post("http://apiserver", obj)
func Post(baseURL string, object *jsh.Object) (*jsh.Document, *http.Response, error) {

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing URL: %s", err.Error())
	}

	// ghetto pluralization, fix when it becomes an issue
	setPath(u, object.Type)

	request, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Error building POST request: %s", err.Error())
	}

	return doObjectRequest(request, object)
}
