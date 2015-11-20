// Package jsc stands for JSON Specification Client. As opposed to the jsh package which is
// namely for parsing incoming requests, this does the opposite in that it is a
// useful wrapper around the base net/http.Client for sending HTTP requests.
package jsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/derekdowling/go-json-spec-handler"
)

// Request is just a wrapper around an http.Request to make sending more fluent
type Request struct {
	request *http.Request
}

// Send sends an http.Request and handles parsing the response back
func (r *Request) Send() (*http.Response, error) {
	client := &http.Client{}
	return client.Do(r.request)
}

// NewRequest creates a new JSON Spec compatible http.Request. Can be used like so:
//
//  obj, err := jsh.NewObject("123", "objtype", payload)
//  req, err := jsc.NewRequest("POST", "http://postap.com", obj)
//  resp, err := req.Send()
//
func NewRequest(method string, urlStr string, object *jsh.Object) (*Request, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch method {
	case "PATCH":
	case "DELETE":
		u.Path = strings.Join([]string{object.Type, object.ID}, "/")
		break
	case "POST":
		break
	default:
		return nil, jsh.SpecificationError(fmt.Sprintf(
			"Cannot use HTTP method ''%s' for a JSON Request", method,
		))
	}

	request, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating new HTTP request: %s", err.Error())
	}

	// use Prepare to generate a payload
	payload, err := object.Prepare(request)
	if err != nil {
		return nil, fmt.Errorf("Error preparing object: %s", err.Error())
	}

	content, jsonErr := json.MarshalIndent(payload, "", "  ")
	if jsonErr != nil {
		return nil, fmt.Errorf("Unable to prepare JSON content: %s", jsonErr)
	}

	request.Body = jsh.CreateReadCloser(content)
	request.Header.Add("Content-Type", jsh.ContentType)
	request.Header.Set("Content-Length", strconv.Itoa(len(content)))

	return &Request{request}, nil
}
