package jsh

import "net/http"

// List is just a wrapper around an object array that implements Sendable
type List struct {
	Objects []*Object `json:"data"`
}

// Prepare returns a success status response
func (l List) Prepare(r *http.Request) (*Response, *Error) {
	return &Response{Data: l.Objects, HTTPStatus: http.StatusOK}, nil
}
