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
	// Status is the HTTP Status Code that should be associated with the object
	// when it is sent.
	Status int `json:"-"`
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

/*
Unmarshal puts an Object's Attributes into a more useful target resourceType defined
by the user. A correct object resourceType specified must also be provided otherwise
an error is returned to prevent hard to track down situations.

Optionally, used https://github.com/go-validator/validator for request input validation.
Simply define your struct with valid input tags:

	struct {
		Username string `json:"username" valid:"required,alphanum"`
	}


As the final action, the Unmarshal function will run govalidator on the unmarshal
result. If the validator fails, a Sendable error response of HTTP Status 422 will
be returned containing each validation error with a populated Error.Source.Pointer
specifying each struct attribute that failed. In this case, all you need to do is:

	errors := obj.Unmarshal("mytype", &myType)
	if errors != nil {
		// log errors via error.ISE
		jsh.Send(w, r, errors)
	}
*/
func (o *Object) Unmarshal(resourceType string, target interface{}) []*Error {

	if resourceType != o.Type {
		return []*Error{ISE(fmt.Sprintf(
			"Expected type %s, when converting actual type: %s",
			resourceType,
			o.Type,
		))}
	}

	jsonErr := json.Unmarshal(o.Attributes, target)
	if jsonErr != nil {
		return []*Error{ISE(fmt.Sprintf(
			"For type '%s' unable to marshal: %s\nError:%s",
			resourceType,
			string(o.Attributes),
			jsonErr.Error(),
		))}
	}

	return validateInput(target)
}

/*
Marshal allows you to load a modified payload back into an object to preserve
all of the data it has.
*/
func (o *Object) Marshal(attributes interface{}) *Error {
	raw, err := json.MarshalIndent(attributes, "", " ")
	if err != nil {
		return ISE(fmt.Sprintf("Error marshaling attrs while creating a new JSON Object: %s", err))
	}

	o.Attributes = raw
	return nil
}

/*
Validate ensures that an object is JSON API compatible. Has a side effect of also
setting the Object's Status attribute to be used as the Response HTTP Code if one
has not already been set.
*/
func (o *Object) Validate(r *http.Request, response bool) *Error {

	if o.ID == "" {

		// don't error if the client is attempting to performing a POST request, in
		// which case, ID shouldn't actually be set
		if !response && r.Method != "POST" {
			return SpecificationError("ID must be set for Object response")
		}
	}

	if o.Type == "" {
		return SpecificationError("Type must be set for Object response")
	}

	switch r.Method {
	case "POST":
		acceptable := map[int]bool{201: true, 202: true, 204: true}

		if o.Status != 0 {
			if _, validCode := acceptable[o.Status]; !validCode {
				return SpecificationError("POST Status must be one of 201, 202, or 204.")
			}
			break
		}

		o.Status = http.StatusCreated
		break
	case "PATCH":
		acceptable := map[int]bool{200: true, 202: true, 204: true}

		if o.Status != 0 {
			if _, validCode := acceptable[o.Status]; !validCode {
				return SpecificationError("PATCH Status must be one of 200, 202, or 204.")
			}
			break
		}

		o.Status = http.StatusOK
		break
	case "GET":
		o.Status = http.StatusOK
		break
	// If we hit this it means someone is attempting to use an unsupported HTTP
	// method. Return a 406 error instead
	default:
		return SpecificationError(fmt.Sprintf(
			"The JSON Specification does not accept '%s' requests.",
			r.Method,
		))
	}

	return nil
}

// String prints a formatted string representation of the object
func (o *Object) String() string {
	raw, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return err.Error()
	}

	return string(raw)
}

// validateInput runs go-validator on each attribute on the struct and returns all
// errors that it picks up
func validateInput(target interface{}) []*Error {

	_, validationError := govalidator.ValidateStruct(target)
	if validationError != nil {

		errorList, isType := validationError.(govalidator.Errors)
		if isType {

			errors := []*Error{}
			for _, singleErr := range errorList.Errors() {

				// parse out validation error
				goValidErr, _ := singleErr.(govalidator.Error)
				inputErr := InputError(goValidErr.Err.Error(), goValidErr.Name)

				errors = append(errors, inputErr)
			}

			return errors
		}
	}

	return nil
}
