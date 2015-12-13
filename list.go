package jsh

import "net/http"

// List is just a wrapper around an object array that implements Sendable
type List []*Object

// Prepare returns a success status response
func (list List) Prepare(r *http.Request, response bool) (*Response, *Error) {
	return &Response{Data: list, HTTPStatus: http.StatusOK}, nil
}
