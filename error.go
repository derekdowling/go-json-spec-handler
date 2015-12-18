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

// ErrorObject consists of a number of contextual attributes to make conveying
// certain error type simpler as per the JSON API specification:
// http://jsonapi.org/format/#error-objects
//
//	error := &jsh.Error{
//		Title: "Authentication Failure",
//		Detail: "Category 4 Username Failure",
//		Status: 401
//	}
//
//	jsh.Send(w, r, error)
//
type ErrorObject struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
	ISE string `json:"-"`
}

// Error is a safe for public consumption error message
func (e *ErrorObject) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Title, e.Detail)
	if e.Source.Pointer != "" {
		msg += fmt.Sprintf("Source.Pointer: %s", e.Source.Pointer)
	}
	return msg
}

// Internal is a convenience function that prints out the full error including the
// ISE which is useful when debugging, NOT to be used for returning errors to user,
// use e.Error() for that
func (e *ErrorObject) Internal() string {
	return fmt.Sprintf("%s ISE: %s", e.Error(), e.ISE)
}

// Error is a Sendable type consistenting of one or more error messages. Error
// implements Sendable and as such, when encountered, can simply be sent via
// jsh:
//
//	object, err := ParseObject(request)
//	if err != nil {
//		err := jsh.Send(err, w, request)
//	}
type Error struct {
	Objects []*ErrorObject
}

// Error allows ErrorList to conform to the default Go error interface
func (e *Error) Error() string {
	err := "Errors: "
	for _, m := range e.Objects {
		err = strings.Join([]string{err, fmt.Sprintf("%s;", m.Error())}, "\n")
	}
	return err
}

// Status returns the HTTP Code of the first Error Object, or 0 if none
func (e *Error) Status() int {
	if len(e.Objects) > 0 {
		return e.Objects[0].Status
	}

	return 0
}

// Internal prints a formatted error list including ISE's, useful for debugging
func (e *Error) Internal() string {
	err := "Errors:"
	for _, m := range e.Objects {
		err = strings.Join([]string{err, fmt.Sprintf("%s;", m.Internal())}, "\n")
	}
	return err
}

// Add first validates the error, and then appends it to the ErrorList
func (e *Error) Add(object *ErrorObject) *Error {
	err := validateError(object)
	if err != nil {
		return err
	}

	e.Objects = append(e.Objects, object)
	return nil
}

// Prepare first validates the errors, and then returns an appropriate response
func (e *Error) Prepare(req *http.Request, response bool) (*Response, *Error) {
	if len(e.Objects) == 0 {
		return nil, ISE("No errors provided for attempted error response.")
	}

	return &Response{Errors: e.Objects, HTTPStatus: e.Objects[0].Status}, nil
}

// validateError ensures that the error is ready for a response in it's current state
func validateError(err *ErrorObject) *Error {

	if err.Status < 400 || err.Status > 600 {
		return ISE(fmt.Sprintf("Invalid HTTP Status for error %+v\n", err))
	} else if err.Status == 422 && err.Source.Pointer == "" {
		return ISE(fmt.Sprintf("Source Pointer must be set for 422 Status errors"))
	}

	return nil
}

// NewError is a convenience function that makes creating a Sendable Error from a
// Error Object simple. Because ErrorObjects are validated agains the JSON API
// Specification before being added, there is a chance that a ISE error might be
// returned in your new error's place.
func NewError(object *ErrorObject) *Error {
	newError := &Error{}

	err := newError.Add(object)
	if err != nil {
		return err
	}

	return newError
}

// ISE is a convenience function for creating a ready-to-go Internal Service Error
// response. The message you pass in is set to the ErrorObject.ISE attribute so you
// can gracefully log ISE's internally before sending them
func ISE(internalMessage string) *Error {
	return NewError(&ErrorObject{
		Title:  DefaultErrorTitle,
		Detail: DefaultErrorDetail,
		Status: http.StatusInternalServerError,
		ISE:    internalMessage,
	})
}

// InputError creates a properly formatted Status 422 error with an appropriate
// user facing message, and a Status Pointer to the first attribute that
func InputError(attribute string, detail string) *Error {
	message := &ErrorObject{
		Title:  "Invalid Attribute",
		Detail: detail,
		Status: 422,
	}

	// Assign this after the fact, easier to do
	message.Source.Pointer = fmt.Sprintf("/data/attributes/%s", strings.ToLower(attribute))

	err := &Error{}
	err.Add(message)
	return err
}

// SpecificationError is used whenever the Client violates the JSON API Spec
func SpecificationError(detail string) *Error {
	return NewError(&ErrorObject{
		Title:  "JSON API Specification Error",
		Detail: detail,
		Status: http.StatusNotAcceptable,
	})
}

// NotFound returns a 404 formatted error
func NotFound(resourceType string, id string) *Error {
	return NewError(&ErrorObject{
		Title:  "Not Found",
		Detail: fmt.Sprintf("No resource of type '%s' exists for ID: %s", resourceType, id),
		Status: http.StatusNotFound,
	})
}
