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
	request, err := PostRequest(baseURL, object)
	if err != nil {
		return nil, nil, err
	}

	return Do(request, jsh.ObjectMode)
}

// PostRequest returns a fully formatted request with JSON body for performing
// a JSONAPI POST. This is useful for if you need to set custom headers on the
// request. Otherwise just use "jsc.Post".
func PostRequest(baseURL string, object *jsh.Object) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL: %s", err.Error())
	}

	setPath(u, object.Type)

	request, err := NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building POST request: %s", err.Error())
	}

	err = prepareBody(request, object)
	if err != nil {
		return nil, err
	}

	return request, nil
}
