package jsh

import "net/http"

// List is just a wrapper around an object array that implements Sendable
type List struct {
	Objects []*Object `json:"data"`
}

// Prepare returns a success status response
func (l *List) prepare(r *http.Request, response bool) (*Response, SendableError) {
	return &Response{Data: l.Objects, HTTPStatus: http.StatusOK}, nil
}

// Add is just a convenience function that appends an additional object to a list
func (l *List) Add(o *Object) {
	l.Objects = append(l.Objects, o)
}
