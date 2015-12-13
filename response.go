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
	// Prepare allows a "raw" response type to perform specification assertions,
	// and format any data before it is actually send
	Prepare(r *http.Request, response bool) (*Response, *Error)
}

// Response represents the top level json format of incoming requests
// and outgoing responses. Refer to the JSON API Specification for a full descriptor
// of each attribute: http://jsonapi.org/format/#document-structure
type Response struct {
	Data     interface{} `json:"data,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	Links    *Link       `json:"links,omitempty"`
	Included *List       `json:"included,omitempty"`
	JSONAPI  struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	// Custom HTTPStatus attribute to make it simpler to define expected responses
	// for a given context
	HTTPStatus int `json:"-"`
	// empty is used to signify that the response shouldn't contain a json payload
	// in the case that we only want to return an HTTP Status Code in order to bypass
	// validation steps
	empty bool
}

// Validate checks JSON Spec for the top level JSON document
func (r *Response) Validate() *Error {

	if !r.empty {
		if r.Errors == nil && r.Data == nil {
			return ISE("Both `errors` and `data` cannot be blank for a JSON response")
		}
		if r.Errors != nil && r.Data != nil {
			return ISE("Both `errors` and `data` cannot be set for a JSON response")
		}
		if r.Data == nil && r.Included != nil {
			return ISE("'included' should only be set for a response if 'data' is as well")
		}
	}
	if r.HTTPStatus < 100 || r.HTTPStatus > 600 {
		return ISE("Response HTTP Status is outside of valid range")
	}

	// probably not the best place for this, but...
	r.JSONAPI.Version = JSONAPIVersion

	return nil
}

// Send will return a JSON payload to the requestor. If the payload response validation
// fails, it will send an appropriate error to the requestor and will return the error
func Send(w http.ResponseWriter, r *http.Request, payload Sendable) *Error {

	response, err := payload.Prepare(r, true)
	if err != nil {

		// use the prepared error as the new response, unless something went horribly
		// wrong
		response, err = err.Prepare(r, true)
		if err != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return err
		}
	}

	return SendResponse(w, r, response)
}

// SendResponse handles sending a fully packaged JSON Response which is useful if you
// require manual, or custom, validation when building a Response:
//
//	response := &jsh.Response{
//		HTTPStatus: http.StatusAccepted,
//	}
//
// The function will always send but will return the last error it encountered
// to help with debugging in the event of an ISE.
func SendResponse(w http.ResponseWriter, r *http.Request, response *Response) *Error {

	err := response.Validate()
	if err != nil {
		errResp, prepErr := err.Prepare(r, true)

		// If we ever hit this, something seriously wrong has happened
		if prepErr != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return prepErr
		}

		// if we didn't error out, make this the new response
		response = errResp
	}

	content, jsonErr := json.MarshalIndent(response, "", "  ")
	if jsonErr != nil {
		http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
		return ISE(fmt.Sprintf("Unable to marshal JSON payload: %s", jsonErr.Error()))
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(response.HTTPStatus)
	w.Write(content)

	if err != nil {
		return err
	}

	return nil
}

// OkResponse fulfills the Sendable interface for a simple success response
type OkResponse struct{}

// Ok makes it simple to return a 200 OK response via jsh:
//
//	jsh.Send(w, r, jsh.Ok())
func Ok() *OkResponse {
	return &OkResponse{}
}

// Prepare turns OkResponse into the normalized Response type
func (o *OkResponse) Prepare(r *http.Request, response bool) (*Response, *Error) {
	return &Response{HTTPStatus: http.StatusOK, empty: true}, nil
}
