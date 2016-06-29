// Package jsh (JSON API Specification Handler) makes it easy to parse JSON API
// requests and send responses that match the JSON API Specification: http://jsonapi.org/
// from your server.
//
// For a request client, see: jsc: https://godoc.org/github.com/derekdowling/go-json-spec-handler/client
//
// For a full http.Handler API builder see jshapi: https://godoc.org/github.com/derekdowling/go-json-spec-handler/jsh-api
package jsh

const (
	// ContentType is the data encoding of choice for HTTP Request and Response Headers
	ContentType = "application/vnd.api+json"
)
