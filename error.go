package jsh

import (
	"fmt"
	"net/http"
	"strings"
)

// DefaultError can be customized in order to provide a more customized error
// Detail message when an Internal Server Error occurs. Optionally, you can modify
// a returned jsh.Error before sending it as a response as well.
var DefaultErrorDetail = "Request failed, something went wrong."

// DefaultTitle can be customized to provide a more customized ISE Title
var DefaultErrorTitle = "Internal Server Error"

// SendableError conforms to a standard error format for logging, but can also
// be sent as a JSON response
type SendableError interface {
	Sendable
	Error() string
}

// Error represents a JSON Specification Error. Error.Source.Pointer is used in 422
// status responses to indicate validation errors on a JSON Object attribute.
//
// ISE (internal server error) captures the server error internally to help with
// logging/troubleshooting, but is never returned in a response.
//
// Once a jsh.Error is returned, and you have logged/handled it accordingly, you
// can simply return it using jsh.Send():
//
//	error := &jsh.Error{
//		Title: "Authentication Failure",
//		Detail: "Category 4 Username Failure",
//		Status: 401
//	}
//
//	jsh.Send(w, r, error)
//
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
	ISE string `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s. %s", e.Title, e.Detail, e.Source.Pointer)
}

// Prepare returns a response containing a prepared error list since the JSON
// API specification requires that errors are returned as a list
func (e *Error) prepare(req *http.Request, response bool) (*Response, SendableError) {
	list := &ErrorList{Errors: []*Error{e}}
	return list.prepare(req, response)
}

// ErrorList is just a wrapped error array that implements Sendable
type ErrorList struct {
	Errors []*Error
}

// Error allows ErrorList to conform to the default Go error interface
func (e *ErrorList) Error() string {
	err := "Errors: "
	for _, e := range e.Errors {
		err = fmt.Sprintf("%s%s;", err, e.Error())
	}
	return err
}

// Add first validates the error, and then appends it to the ErrorList
func (e *ErrorList) Add(newError *Error) *Error {
	err := validateError(newError)
	if err != nil {
		return err
	}

	e.Errors = append(e.Errors, newError)
	return nil
}

// Prepare first validates the errors, and then returns an appropriate response
func (e *ErrorList) prepare(req *http.Request, response bool) (*Response, SendableError) {
	if len(e.Errors) == 0 {
		return nil, ISE("No errors provided for attempted error response.")
	}

	return &Response{Errors: e.Errors, HTTPStatus: e.Errors[0].Status}, nil
}

// validateError ensures that the error is ready for a response in it's current state
func validateError(err *Error) *Error {

	if err.Status < 400 || err.Status > 600 {
		return ISE(fmt.Sprintf("Invalid HTTP Status for error %+v\n", err))
	} else if err.Status == 422 && err.Source.Pointer == "" {
		return ISE(fmt.Sprintf("Source Pointer must be set for 422 Status errors"))
	}

	return nil
}

// ISE is a convenience function for creating a ready-to-go Internal Service Error
// response. As previously mentioned, the Error.ISE field is for logging only, and
// won't be returned to the end user.
func ISE(err string) *Error {
	return &Error{
		Title:  DefaultErrorTitle,
		Detail: DefaultErrorDetail,
		Status: http.StatusInternalServerError,
		ISE:    err,
	}
}

// InputError creates a properly formatted Status 422 error with an appropriate
// user facing message, and a Status Pointer to the first attribute that
func InputError(attribute string, detail string) *Error {
	err := &Error{
		Title:  "Invalid Attribute",
		Detail: detail,
		Status: 422,
	}

	// Assign this after the fact, easier to do
	err.Source.Pointer = fmt.Sprintf("/data/attributes/%s", strings.ToLower(attribute))

	return err
}

// SpecificationError is used whenever the Client violates the JSON API Spec
func SpecificationError(detail string) *Error {
	return &Error{
		Title:  "API Specification Error",
		Detail: detail,
		Status: http.StatusNotAcceptable,
	}
}
