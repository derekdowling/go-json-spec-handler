package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// DocumentMode allows different specification settings to be enforced
// based on the specified mode.
type DocumentMode int

const (
	// ObjectMode enforces fetch request/response specifications
	ObjectMode DocumentMode = iota
	// ListMode enforces listing request/response specifications
	ListMode
	// ErrorMode enforces error response specifications
	ErrorMode
)

// IncludeJSONAPIVersion is an option that allows consumers to include/remove the `jsonapi`
// top-level member from server responses.
var IncludeJSONAPIVersion = true

// JSONAPI is the top-level member of a JSONAPI document that includes
// the server compatible version of the JSONAPI specification.
type JSONAPI struct {
	Version string `json:"version"`
}

/*
Document represents a top level JSON formatted Document.
Refer to the JSON API Specification for a full descriptor
of each attribute: http://jsonapi.org/format/#document-structure
*/
type Document struct {
	Data List `json:"data"`
	// Object   *Object     `json:"-"`
	Errors   ErrorList   `json:"errors,omitempty"`
	Links    *Links      `json:"links,omitempty"`
	Included []*Object   `json:"included,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	JSONAPI  *JSONAPI    `json:"jsonapi,omitempty"`
	// Status is an HTTP Status Code
	Status int `json:"-"`
	// DataMode to enforce for the document
	Mode DocumentMode `json:"-"`
	// empty is used to signify that the response shouldn't contain a json payload
	// in the case that we only want to return an HTTP Status Code in order to bypass
	// validation steps.
	empty bool
	// validated confirms whether or not the document as a whole is validated and
	// in a safe-to-send state.
	validated bool
}

/*
New instantiates a new JSON Document object.
*/
func New() *Document {
	json := &Document{}
	if IncludeJSONAPIVersion {
		json.JSONAPI = &JSONAPI{
			Version: JSONAPIVersion,
		}
	}

	return json
}

/*
Build creates a Sendable Document with the provided sendable payload, either Data or
errors. Build also assumes you've already validated your data with .Validate() so
it should be used carefully.
*/
func Build(payload Sendable) *Document {
	document := New()
	document.validated = true

	object, isObject := payload.(*Object)
	if isObject {
		document.Data = List{object}
		document.Status = object.Status
		document.Mode = ObjectMode
	}

	list, isList := payload.(List)
	if isList {
		document.Data = list
		document.Status = http.StatusOK
		document.Mode = ListMode
	}

	err, isError := payload.(*Error)
	if isError {
		document.Errors = ErrorList{err}
		document.Status = err.Status
		document.Mode = ErrorMode
	}

	errorList, isErrorList := payload.(ErrorList)
	if isErrorList {
		document.Errors = errorList
		document.Status = errorList[0].Status
		document.Mode = ErrorMode
	}

	return document
}

/*
Validate performs document level checks against the JSONAPI specification. It is
assumed that if this call returns without an error, your document is valid and
can be sent as a request or response.
*/
func (d *Document) Validate(r *http.Request, isResponse bool) *Error {

	// if sending a response, we must have a valid HTTP status at the very least
	// to send
	if isResponse && d.Status < 100 || d.Status > 600 {
		return ISE("Response HTTP Status is outside of valid range")
	}

	// There are certain cases such as HTTP 204 that send without a payload,
	// this is the short circuit to make sure we don't false alarm on those cases
	if d.empty {
		return nil
	}

	// if we have errors, and they have been added in a way that does not trigger
	// error mode, set it now so we perform the proper validations.
	if d.HasErrors() && d.Mode != ErrorMode {
		d.Mode = ErrorMode
	}

	switch d.Mode {
	case ErrorMode:
		if d.HasData() {
			return ISE("Attempting to respond with 'data' in an error response")
		}
	case ObjectMode:
		if d.HasData() && len(d.Data) > 1 {
			return ISE("Cannot set more than one data object in 'ObjectMode'")
		}
	case ListMode:
		if !d.HasErrors() && d.Data == nil {
			return ISE("Data cannot be nil in 'ListMode', use empty array")
		}
	}

	if !d.HasData() && d.Included != nil {
		return ISE("'included' should only be set for a response if 'data' is as well")
	}

	err := d.Data.Validate(r, isResponse)
	if err != nil {
		return err
	}

	err = d.Errors.Validate(r, isResponse)
	if err != nil {
		return err
	}

	d.validated = true

	return nil
}

// AddObject adds another object to the JSON Document after validating it.
func (d *Document) AddObject(object *Object) *Error {

	if d.Mode == ErrorMode {
		return ISE("Cannot add data to a document already possessing errors")
	}

	if d.Mode == ObjectMode && len(d.Data) == 1 {
		return ISE("Single 'data' object response is expected, you are attempting to add more than one element to be returned")
	}

	// if not yet set, add the associated HTTP status with the object
	if d.Status == 0 {
		d.Status = object.Status
	}

	// finally, actually add the object to data List
	if d.Data == nil {
		d.Data = List{object}
	} else {
		d.Data = append(d.Data, object)
	}

	return nil
}

/*
AddError adds an error to the Document. It will also set the document Mode to
"ErrorMode" if not done so already.
*/
func (d *Document) AddError(newErr *Error) *Error {

	if d.HasData() {
		return ISE("Attempting to set an error, when the document has prepared response data")
	}

	if newErr.Status == 0 {
		return ISE("No HTTP Status code provided for error, cannot add to document")
	}

	if d.Status == 0 {
		d.Status = newErr.Status
	}

	if d.Errors == nil {
		d.Errors = []*Error{newErr}
	} else {
		d.Errors = append(d.Errors, newErr)
	}

	// set document to error mode
	d.Mode = ErrorMode

	return nil
}

/*
First will return the first object from the document data if possible.
*/
func (d *Document) First() *Object {
	if !d.HasData() {
		return nil
	}

	return d.Data[0]
}

// HasData will return true if the JSON document's Data field is set
func (d *Document) HasData() bool {
	return d.Data != nil && len(d.Data) > 0
}

// HasErrors will return true if the Errors attribute is not nil.
func (d *Document) HasErrors() bool {
	return d.Errors != nil && len(d.Errors) > 0
}

func (d *Document) Error() string {
	errStr := "Errors:"
	for _, err := range d.Errors {
		errStr = strings.Join([]string{errStr, fmt.Sprintf("%s;", err.Error())}, "\n")
	}
	return errStr
}

/*
MarshalJSON handles the custom serialization case caused by case where the "data"
element of a document might be either a single resource object, or a collection of
them.
*/
func (d *Document) MarshalJSON() ([]byte, error) {
	// we use the MarshalDoc type to avoid recursively calling this function below
	// when we marshal
	type MarshalDoc Document
	doc := MarshalDoc(*d)

	switch d.Mode {
	case ObjectMode:
		var data *Object
		if len(d.Data) > 0 {
			data = d.Data[0]
		}

		// subtype that overrides regular data List with a single Object for
		// fetch style request/responses
		type MarshalObject struct {
			MarshalDoc
			Data *Object `json:"data"`
		}

		return json.Marshal(MarshalObject{
			MarshalDoc: doc,
			Data:       data,
		})

	case ErrorMode:
		// subtype that omits data as expected for error responses. We cannot simply
		// use json:"-" for the data attribute otherwise it will not override the
		// default struct tag of it the composed MarshalDoc struct.
		type MarshalError struct {
			MarshalDoc
			Data *Object `json:"data,omitempty"`
		}

		return json.Marshal(MarshalError{
			MarshalDoc: doc,
		})

	case ListMode:
		return json.Marshal(doc)
	default:
		return nil, ISE(fmt.Sprintf("Unexpected DocumentMode value when marshaling: %d", d.Mode))
	}
}
