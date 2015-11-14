package jsh

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/validator.v2"
)

// Object represents the default JSON spec for objects
type Object struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    json.RawMessage    `json:"attributes,omitempty"`
	Links         map[string]*Link   `json:"links,omitempty"`
	Relationships map[string]*Object `json:"relationships,omitempty"`
}

// NewObject prepares a new JSON Object for an API response. Whatever is provided
// as attributes will be marshalled to JSON.
func NewObject(id string, objType string, attributes interface{}) (*Object, error) {
	object := &Object{
		ID:            id,
		Type:          objType,
		Links:         map[string]*Link{},
		Relationships: map[string]*Object{},
	}

	rawJSON, err := json.MarshalIndent(attributes, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Error marshaling attrs while creating a new JSON Object: %s", err)
	}

	object.Attributes = rawJSON
	return object, nil
}

// Unmarshal puts an Object's Attributes into a more useful target type defined
// by the user. A correct object type specified must also be provided otherwise
// an error is returned to prevent hard to track down situations.
//
// Optionally, used https://github.com/go-validator/validator for request input validation.
// Simply define your struct with valid input tags:
//
//		struct {
//			Username string `validate:"min=3,max=40,regexp=^[a-zA-Z]$"`
//		}
//
//	and the function will run go-validator on the unmarshal result. If the validator
//	fails, a Sendable error response of HTTP Status 422 will be returned containing
//	each validation error with a populated Error.Source.Pointer specifying each struct
//	attribute that failed. In this case, all you need to do is:
//
//		errors := obj.Unmarshal("mytype", &myType)
//		if errors != nil {
//			// log errors via error.ISE
//			jsh.Send(r, w, errors)
//		}
func (o *Object) Unmarshal(objType string, target interface{}) (err Sendable) {

	if objType != o.Type {
		err = ISE(fmt.Sprintf(
			"Expected type %s, when converting actual type: %s",
			objType,
			o.Type,
		))
		return
	}

	jsonErr := json.Unmarshal(o.Attributes, target)
	if jsonErr != nil {
		err = ISE(fmt.Sprintf(
			"For type '%s' unable to marshal: %s\nError:%s",
			objType,
			string(o.Attributes),
			jsonErr.Error(),
		))
		return
	}

	ok, errors := validateInput(target)
	if !ok {
		return errors
	}

	return nil
}

// Prepare creates a new JSON single object response with an appropriate HTTP status
// to match the request method type.
func (o *Object) Prepare(r *http.Request) (*Response, *Error) {

	var status int

	switch r.Method {
	case "POST":
		status = http.StatusCreated
	case "PATCH":
		status = http.StatusOK
	case "GET":
		status = http.StatusOK
	// If we hit this it means someone is attempting to use an unsupported HTTP
	// method. Return a 406 error instead
	default:
		return SpecificationError(fmt.Sprintf(
			"The JSON Specification does not accept '%s' requests.",
			r.Method,
		)).Prepare(r)
	}

	return &Response{HTTPStatus: status, Data: o}, nil
}

// validateInput runs go-validator on each attribute on the struct and returns all
// errors that it picks up
func validateInput(target interface{}) (ok bool, errors *ErrorList) {
	ok = true
	errors = &ErrorList{}

	err := validator.Validate(target)
	if err != nil {
		ok = false
		log.Printf("errors = %+v\n", errors)

		// Each attribute can have multiple errors, only return the first one for each
		// for attributeName, attributeErrors := range errs {
		// attributeError := attributeErrors[0]
		// errors.Add(InputError(attributeName, attributeError))
		// }
	}

	return
}
