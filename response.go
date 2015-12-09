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
	Prepare(r *http.Request, response bool) (*Response, SendableError)
}

// Response represents the top level json format of incoming requests
// and outgoing responses
type Response struct {
	HTTPStatus int         `json:"-"`
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
func (r *Response) Validate() SendableError {

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

// Send will return a JSON payload to the requestor. If the payload response validation
// fails, it will send an appropriate error to the requestor and will return the error
func Send(w http.ResponseWriter, r *http.Request, payload Sendable) error {

	response, err := payload.Prepare(r, true)
	if err != nil {

		response, err = err.Prepare(r, true)
		if err != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return fmt.Errorf("Error preparing JSH error: %s", err.Error())
		}

		return fmt.Errorf("Error preparing JSON payload: %s", err.Error())
	}

	return SendResponse(w, r, response)
}

// SendResponse handles sending a fully packaged JSON Response allows API consumers
// to more manually build their Responses in case they want to send Meta, Links, etc
// The function will always, send but will return the last error it encountered
// to help with debugging
func SendResponse(w http.ResponseWriter, r *http.Request, response *Response) error {

	err := response.Validate()
	if err != nil {
		response, err = err.Prepare(r, true)

		// If we ever hit this, something seriously wrong has happened
		if err != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return fmt.Errorf("Error preparing JSH error: %s", err.Error())
		}

		return fmt.Errorf("Response validation error: %s", err.Error())
	}

	content, jsonErr := json.MarshalIndent(response, "", "  ")
	if jsonErr != nil {
		http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
		return fmt.Errorf("Unable to marshal JSON payload: %s", jsonErr.Error())
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(response.HTTPStatus)
	w.Write(content)

	return nil
}
