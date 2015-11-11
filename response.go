package japi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse for API requests
type ErrorResponse struct {
	Errors []*Error `json:"errors"`
}

// RespondWithObject sends a single data object as a JSON response
func RespondWithObject(w http.ResponseWriter, status int, object *Object) error {
	payload := struct {
		Data *Object `json:"data"`
	}{object}

	return respond(w, status, payload)
}

// RespondWithList sends a list of data objects as a JSON response
func RespondWithList(w http.ResponseWriter, status int, list []*Object) error {
	payload := struct {
		Data []*Object `json:"data"`
	}{list}
	return respond(w, status, payload)
}

// RespondWithError is a convenience function that puts an error into an array
// and then calls RespondWithErrors which is the correct error response format
func RespondWithError(w http.ResponseWriter, err *Error) error {
	errors := []*Error{err}
	return RespondWithErrors(w, errors)
}

// RespondWithErrors sends the expected error response format for a
// request that cannot be fulfilled in someway. Allows the user
// to compile multiple errors that can be sent back to a user. Uses
// the first error status as the HTTP Status to return.
func RespondWithErrors(w http.ResponseWriter, errors []*Error) error {

	if len(errors) == 0 {
		return fmt.Errorf("No errors provided for attempted error response.")
	}

	// use the first error to set the error status
	status := errors[0].Status
	payload := ErrorResponse{errors}
	return respond(w, status, payload)
}

// Respond formats a JSON response with the appropriate headers.
func respond(w http.ResponseWriter, status int, payload interface{}) error {
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

func prepareError(error *Error) *ErrorResponse {
	return &ErrorResponse{}
}

func prepareObject(object *Object) *DataResponse {
	return &DataResponse{object}
}

func prepareList(list []*Object) *DataResponse {
	return &DataResponse{list}
}
