package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
Document represents a top level JSON formatted Document.
Refer to the JSON API Specification for a full descriptor
of each attribute: http://jsonapi.org/format/#document-structure
*/
type Document struct {
	Data     List        `json:"data"`
	Object   *Object     `json:"-"`
	Errors   ErrorList   `json:"errors,omitempty"`
	Links    *Link       `json:"links,omitempty"`
	Included []*Object   `json:"included,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	JSONAPI  struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	// Status is an HTTP Status Code
	Status int `json:"-"`
	// empty is used to signify that the response shouldn't contain a json payload
	// in the case that we only want to return an HTTP Status Code in order to bypass
	// validation steps.
	empty     bool
	validated bool
}

/*
New instantiates a new JSON Document object.
*/
func New() *Document {
	json := &Document{}
	json.JSONAPI.Version = JSONAPIVersion

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
		document.Object = object
		document.Status = object.Status
	}

	list, isList := payload.(List)
	if isList {
		document.Data = list
		document.Status = http.StatusOK
	}

	err, isError := payload.(*Error)
	if isError {
		document.Errors = ErrorList{err}
		document.Status = err.Status
	}

	errorList, isErrorList := payload.(ErrorList)
	if isErrorList {
		document.Errors = errorList
		document.Status = errorList[0].Status
	}

	return document
}

/*
Validate checks JSON Spec for the top level JSON document
*/
func (d *Document) Validate(r *http.Request, response bool) *Error {

	if d.Status < 100 || d.Status > 600 {
		return ISE("Response HTTP Status is outside of valid range")
	}

	// if empty is set, skip all validations below
	if d.empty {
		return nil
	}

	if d.Object != nil && d.Data != nil {
		return ISE("Both `Data` and `Object` cannot be set for a JSON response")
	}
	if !d.HasErrors() && d.Object == nil && d.Data == nil {
		return ISE("Both `errors` and `data` cannot be blank for a JSON response")
	}
	if d.HasErrors() && d.Data != nil {
		return ISE("Both `errors` and `data` cannot be set for a JSON response")
	}
	if !d.HasData() && d.Included != nil {
		return ISE("'included' should only be set for a response if 'data' is as well")
	}

	// if fields have already been validated, skip this part
	if d.validated {
		return nil
	}

	if d.Object != nil {
		if err := d.Object.Validate(r, response); err != nil {
			return err
		}
	}
	err := d.Data.Validate(r, response)
	if err != nil {
		return err
	}

	err = d.Errors.Validate(r, response)
	if err != nil {
		return err
	}

	return nil
}

// AddObject adds another object to the JSON Document after validating it.
func (d *Document) AddObject(object *Object) *Error {

	if d.HasErrors() {
		return ISE("Cannot add data to a document already possessing errors")
	}

	if d.Object != nil {
		return ISE("Cannot add data to a non-collection")
	}

	if d.Status == 0 {
		d.Status = object.Status
	}

	if d.Data == nil {
		d.Data = List{object}
	} else {
		d.Data = append(d.Data, object)
	}

	return nil
}

// AddError adds an error to the JSON Object by transfering it's Error objects.
func (d *Document) AddError(newErr *Error) *Error {

	if d.HasData() {
		return ISE("Cannot add an error to a document already possessing data")
	}

	if newErr.Status == 0 {
		return SpecificationError("Status code must be set for an error")
	}

	if d.Status == 0 {
		d.Status = newErr.Status
	}

	if d.Errors == nil {
		d.Errors = []*Error{newErr}
	} else {
		d.Errors = append(d.Errors, newErr)
	}

	return nil
}

/*
First is just a convenience function that returns the first data object from the
array
*/
func (d *Document) First() *Object {
	if d.Object != nil {
		return d.Object
	}
	return d.Data[0]
}

// HasData will return true if the JSON document's Data field is set
func (d *Document) HasData() bool {
	return d.Object != nil || (d.Data != nil && len(d.Data) > 0)
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
MarshalJSON handles custom serialization required because the "data" element of a document
might be a collection of objects or a single object.
*/
func (d *Document) MarshalJSON() ([]byte, error) {
	// subtype that overrides data with a single object
	type MarshalObject struct {
		Document
		Data *Object `json:"data,omitempty"`
	}

	if d == nil {
		return nil, nil
	}
	if d.Object != nil {
		return json.Marshal(MarshalObject{
			Document: *d,
			Data:     d.Object,
		})
	}

	// avoid stack overflow by using this subtype for marshaling
	type MarshalList Document
	return json.Marshal(MarshalList(*d))
}
