package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
)

// Object represents the default JSON spec for objects
type Object struct {
	Type          string             `json:"type" valid:"alpha,required"`
	ID            string             `json:"id"`
	Attributes    json.RawMessage    `json:"attributes,omitempty"`
	Links         map[string]*Link   `json:"links,omitempty"`
	Relationships map[string]*Object `json:"relationships,omitempty"`
}

// NewObject prepares a new JSON Object for an API response. Whatever is provided
// as attributes will be marshalled to JSON.
func NewObject(id string, resourceType string, attributes interface{}) (*Object, *Error) {
	object := &Object{
		ID:            id,
		Type:          resourceType,
		Links:         map[string]*Link{},
		Relationships: map[string]*Object{},
	}

	rawJSON, err := json.MarshalIndent(attributes, "", " ")
	if err != nil {
		return nil, ISE(fmt.Sprintf("Error marshaling attrs while creating a new JSON Object: %s", err))
	}

	object.Attributes = rawJSON
	return object, nil
}

// Unmarshal puts an Object's Attributes into a more useful target resourceType defined
// by the user. A correct object resourceType specified must also be provided otherwise
// an error is returned to prevent hard to track down situations.
//
// Optionally, used https://github.com/go-validator/validator for request input validation.
// Simply define your struct with valid input tags:
//
//	struct {
//		Username string `json:"username" valid:"required,alphanum"`
//	}
//
//
// As the final action, the Unmarshal function will run govalidator on the unmarshal
// result. If the validator fails, a Sendable error response of HTTP Status 422 will
// be returned containing each validation error with a populated Error.Source.Pointer
// specifying each struct attribute that failed. In this case, all you need to do is:
//
//	errors := obj.Unmarshal("mytype", &myType)
//	if errors != nil {
//		// log errors via error.ISE
//		jsh.Send(w, r, errors)
//	}
func (o *Object) Unmarshal(resourceType string, target interface{}) *Error {

	if resourceType != o.Type {
		return ISE(fmt.Sprintf(
			"Expected type %s, when converting actual type: %s",
			resourceType,
			o.Type,
		))
	}

	jsonErr := json.Unmarshal(o.Attributes, target)
	if jsonErr != nil {
		return ISE(fmt.Sprintf(
			"For type '%s' unable to marshal: %s\nError:%s",
			resourceType,
			string(o.Attributes),
			jsonErr.Error(),
		))
	}

	return validateInput(target)
}

// Marshal allows you to load a modified payload back into an object to preserve
// all of the data it has
func (o *Object) Marshal(attributes interface{}) *Error {
	raw, err := json.MarshalIndent(attributes, "", " ")
	if err != nil {
		return ISE(fmt.Sprintf("Error marshaling attrs while creating a new JSON Object: %s", err))
	}

	o.Attributes = raw
	return nil
}

// Prepare creates a new JSON single object response with an appropriate HTTP status
// to match the request method type.
func (o *Object) Prepare(r *http.Request, response bool) (*Response, *Error) {

	if o.ID == "" {

		// don't error if the client is attempting to performing a POST request, in
		// which case, ID shouldn't actually be set
		if !response && r.Method != "POST" {
			return nil, SpecificationError("ID must be set for Object response")
		}
	}

	if o.Type == "" {
		return nil, SpecificationError("Type must be set for Object response")
	}

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
		)).Prepare(r, response)
	}

	return &Response{HTTPStatus: status, Data: o}, nil
}

// validateInput runs go-validator on each attribute on the struct and returns all
// errors that it picks up
func validateInput(target interface{}) *Error {

	_, validationError := govalidator.ValidateStruct(target)
	if validationError != nil {

		errorList, isType := validationError.(govalidator.Errors)
		if isType {

			err := &Error{}
			for _, singleErr := range errorList.Errors() {

				// parse out validation error
				goValidErr, _ := singleErr.(govalidator.Error)
				inputErr := InputError(goValidErr.Name, goValidErr.Err.Error())

				// gross way to do this, but will require a refactor probably
				// to achieve something more elegant
				err.Add(inputErr.Objects[0])
			}

			return err
		}
	}

	return nil
}
