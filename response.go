package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// JSONAPIVersion is version of JSON API Spec that is currently compatible:
// http://jsonapi.org/format/1.1/
const JSONAPIVersion = "1.1"

// Sendable implements functions that allows different response types
// to produce a sendable JSON Response format
type Sendable interface {
	Prepare(r *http.Request) (*Response, SendableError)
}

// Response represents the top level json format of incoming requests
// and outgoing responses
type Response struct {
	HTTPStatus int         `json:"-"`
	Validated  bool        `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Links      *Link       `json:"links,omitempty"`
	Included   *List       `json:"included,omitempty"`
	JSONAPI    struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
}

// Validate checks JSON Spec for the top level JSON document
func (r *Response) Validate() *Error {

	if r.Errors == nil && r.Data == nil {
		return ISE("Both `errors` and `data` cannot be blank for a JSON response")
	}
	if r.Errors != nil && r.Data != nil {
		return ISE("Both `errors` and `data` cannot be set for a JSON response")
	}
	if r.Data == nil && r.Included != nil {
		return ISE("'included' should only be set for a response if 'data' is as well")
	}
	if r.HTTPStatus < 100 || r.HTTPStatus > 600 {
		return ISE("Response HTTP Status must be of a valid range")
	}

	// probably not the best place for this, but...
	r.JSONAPI.Version = JSONAPIVersion

	return nil
}

// Send fires a JSON response if the payload is prepared successfully, otherwise it
// returns an Error which can also be sent.
func Send(r *http.Request, w http.ResponseWriter, payload Sendable) SendableError {
	response, err := payload.Prepare(r)
	if err != nil {
		return err
	}

	return SendResponse(r, w, response)
}

// SendResponse handles sending a fully packaged JSON Response allows API consumers
// to more manually build their Responses in case they want to send Meta, Links, etc
func SendResponse(r *http.Request, w http.ResponseWriter, response *Response) SendableError {

	err := response.Validate()
	if err != nil {
		return err
	}

	content, jsonErr := json.MarshalIndent(response, "", "  ")
	if jsonErr != nil {
		// Sendception
		return ISE(fmt.Sprintf("Unable to prepare payload JSON: %s", jsonErr))
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(response.HTTPStatus)
	w.Write(content)

	return nil
}
