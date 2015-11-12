package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Data represents the top level json format of incoming requests
// and outgoing responses
type Data struct {
	Data interface{} `json:"data"`
}

// ErrorResponse for API requests
type ErrorResponse struct {
	Errors []*Error `json:"errors"`
}

// SendObject sends a single data object as a JSON response
func SendObject(w http.ResponseWriter, r *http.Request, object *Object) error {

	var status int

	switch r.Method {
	case "POST":
		status = http.StatusCreated
	case "PATCH":
		status = http.StatusOK
	case "GET":
		status = http.StatusOK
	default:
		return fmt.Errorf("SendObject does not support HTTP Request type: %s", r.Method)
	}

	return Send(w, status, prepareObject(object))
}

// SendList sends a list of data objects as a JSON response
func SendList(w http.ResponseWriter, list []*Object) error {
	return Send(w, http.StatusOK, prepareList(list))
}

// SendError is a convenience function that puts an error into an array
// and then calls SendErrors which is the correct error response format
func SendError(w http.ResponseWriter, err *Error) error {
	return SendErrors(w, prepareError(err))
}

// SendErrors sends the expected error response format for a
// request that cannot be fulfilled in someway. Allows the user
// to compile multiple errors that can be sent back to a user. Uses
// the first error status as the HTTP Status to return.
func SendErrors(w http.ResponseWriter, errors []*Error) error {

	if len(errors) == 0 {
		return fmt.Errorf("No errors provided for attempted error response.")
	}

	for _, err := range errors {
		if err.Status < 400 || err.Status > 500 {
			return fmt.Errorf("Invalid Status for error %+v\n", err)
		} else if err.Status == 422 && err.Source.Pointer == "" {
			return fmt.Errorf("Source Pointer must be set for 422 Status errors")
		}
	}

	// use the first error to set the error status
	status := errors[0].Status
	return Send(w, status, prepareErrorList(errors))
}

// Send formats a JSON response with the appropriate headers.
func Send(w http.ResponseWriter, status int, payload interface{}) error {
	content, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(status)
	w.Write(content)

	return nil
}

func prepareError(err *Error) []*Error {
	return []*Error{err}
}

func prepareErrorList(errors []*Error) *ErrorResponse {
	return &ErrorResponse{Errors: errors}
}

func prepareObject(object *Object) *Data {
	return &Data{Data: object}
}

func prepareList(list []*Object) *Data {
	return &Data{Data: list}
}