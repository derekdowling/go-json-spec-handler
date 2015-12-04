package jsh

import (
	"log"
	"net/http"
)

// List is just a wrapper around an object array that implements Sendable
type List []*Object

// Prepare returns a success status response
func (list List) Prepare(r *http.Request, response bool) (*Response, SendableError) {
	log.Printf("prepare = %+v\n", list)
	return &Response{Data: list, HTTPStatus: http.StatusOK}, nil
}
