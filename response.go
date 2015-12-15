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
	Prepare(r *http.Request, response bool) (*JSON, *Error)
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

	return SendJSON(w, r, response)
}

// SendJSON handles sending a fully prepared JSON Document. This is useful if you
// require custom validation or additional build steps before sending.
//
// SendJSON is designed to always send a response, but will also return the last
// error it encountered to help with debugging in the event of an Internal Server
// Error.
func SendJSON(w http.ResponseWriter, r *http.Request, response *JSON) *Error {

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

	content, jsonErr := json.MarshalIndent(response, "", " ")
	if jsonErr != nil {
		http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
		return ISE(fmt.Sprintf("Unable to marshal JSON payload: %s", jsonErr.Error()))
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(response.HTTPStatus)
	w.Write(content)

	return err
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
func (o *OkResponse) Prepare(r *http.Request, response bool) (*JSON, *Error) {
	return &JSON{HTTPStatus: http.StatusOK, empty: true}, nil
}
