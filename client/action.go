package jsc

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/derekdowling/go-json-spec-handler"
)

// Action performs an outbound GET /resource/:id/action request
func Action(baseURL string, resourceType string, id string, action string) (*jsh.Document, *http.Response, error) {
	request, err := ActionRequest(baseURL, resourceType, id, action)
	if err != nil {
		return nil, nil, err
	}

	return Do(request)
}

/*
ActionRequest returns a fully formatted JSONAPI Action (GET /resource/:id/action) request.
Useful if you need to set custom headers before proceeding. Otherwise just use "jsh.Action".
*/
func ActionRequest(baseURL string, resourceType, id string, action string) (*http.Request, error) {
	if id == "" {
		return nil, jsh.SpecificationError("ID cannot be empty for an Action request type")
	}

	if action == "" {
		return nil, jsh.SpecificationError("Action specifier cannot be empty for an Action request type")
	}

	u, urlErr := url.Parse(baseURL)
	if urlErr != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error parsing URL: %s", urlErr.Error()))
	}

	setIDPath(u, resourceType, id)

	// concat the action the end of url.Path, ensure no "/" prefix
	if strings.HasPrefix(action, "/") {
		action = action[1:]
	}

	u.Path = strings.Join([]string{u.Path, action}, "/")

	return NewRequest("GET", u.String(), nil)
}
