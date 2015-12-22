package jsh

import (
	"fmt"
	"net/http"
	"strings"
)

/*
DefaultError can be customized in order to provide a more customized error
Detail message when an Internal Server Error occurs. Optionally, you can modify
a returned jsh.Error before sending it as a response as well.
*/
var DefaultErrorDetail = "Request failed, something went wrong."

// DefaultTitle can be customized to provide a more customized ISE Title
var DefaultErrorTitle = "Internal Server Error"

/*
Error consists of a number of contextual attributes to make conveying
certain error type simpler as per the JSON API specification:
http://jsonapi.org/format/#error-objects

	error := &jsh.Error{
		Title: "Authentication Failure",
		Detail: "Category 4 Username Failure",
		Status: 401
	}

	jsh.Send(w, r, error)
*/
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
	ISE string `json:"-"`
}

/*
Error will print an internal server error if set, or default back to the SafeError()
format if not. As usual, err.Error() should not be considered safe for presentation
to the end user, use err.SafeError() instead.
*/
func (e *Error) Error() string {
	if e.ISE != "" {
		return fmt.Sprintf("HTTP %d - %s", e.Status, e.ISE)
	}

	return e.SafeError()
}

/*
SafeError is a formatted error string that does not include an associated internal
server error such that it is safe to return this error message to an end user.
*/
func (e *Error) SafeError() string {
	msg := fmt.Sprintf("HTTP %d: %s - %s", e.Status, e.Title, e.Detail)
	if e.Source.Pointer != "" {
		msg += fmt.Sprintf("Source.Pointer: %s", e.Source.Pointer)
	}

	return msg
}

/*
Validate ensures that the an error meets all JSON API criteria.
*/
func (e *Error) Validate(r *http.Request, response bool) *Error {

	if e.Status < 400 || e.Status > 600 {
		return ISE(fmt.Sprintf("Invalid HTTP Status for error %+v\n", e))
	} else if e.Status == 422 && e.Source.Pointer == "" {
		return ISE(fmt.Sprintf("Source Pointer must be set for 422 Status errors"))
	}

	return nil
}

/*
ISE is a convenience function for creating a ready-to-go Internal Service Error
response. The message you pass in is set to the ErrorObject.ISE attribute so you
can gracefully log ISE's internally before sending them.
*/
func ISE(internalMessage string) *Error {
	return &Error{
		Title:  DefaultErrorTitle,
		Detail: DefaultErrorDetail,
		Status: http.StatusInternalServerError,
		ISE:    internalMessage,
	}
}

/*
InputError creates a properly formatted HTTP Status 422 error with an appropriate
user safe message. The parameter "attribute" will format err.Source.Pointer to be
"/data/attributes/<attribute>".
*/
func InputError(msg string, attribute string) *Error {
	err := &Error{
		Title:  "Invalid Attribute",
		Detail: msg,
		Status: 422,
	}

	// Assign this after the fact, easier to do
	err.Source.Pointer = fmt.Sprintf("/data/attributes/%s", strings.ToLower(attribute))

	return err
}

// SpecificationError is used whenever the Client violates the JSON API Spec
func SpecificationError(detail string) *Error {
	return &Error{
		Title:  "JSON API Specification Error",
		Detail: detail,
		Status: http.StatusNotAcceptable,
	}
}

// NotFound returns a 404 formatted error
func NotFound(resourceType string, id string) *Error {
	return &Error{
		Title:  "Not Found",
		Detail: fmt.Sprintf("No resource of type '%s' exists for ID: %s", resourceType, id),
		Status: http.StatusNotFound,
	}
}
