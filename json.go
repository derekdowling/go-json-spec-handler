package jsh

import (
	"encoding/json"
	"log"
)

/*
JSON represents a top level JSON formatted Document.
Refer to the JSON API Specification for a full descriptor
of each attribute: http://jsonapi.org/format/#document-structure
*/
type JSON struct {
	Data     *Data       `json:"data,omitempty"`
	Errors   *Error      `json:"errors,omitempty"`
	Links    *Link       `json:"links,omitempty"`
	Included []*Object   `json:"included,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	JSONAPI  struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	// HTTPStatus attribute to make it simpler to define expected responses
	// for a given context
	HTTPStatus int `json:"-"`
	// empty is used to signify that the response shouldn't contain a json payload
	// in the case that we only want to return an HTTP Status Code in order to bypass
	// validation steps
	empty bool
}

/*
Validate checks JSON Spec for the top level JSON document
*/
func (j *JSON) Validate() *Error {

	if !j.empty {
		if !j.HasErrors() && !j.HasData() {
			return ISE("Both `errors` and `data` cannot be blank for a JSON response")
		}
		if j.HasErrors() && j.HasData() {
			return ISE("Both `errors` and `data` cannot be set for a JSON response")
		}
		if j.HasData() && j.Included != nil {
			return ISE("'included' should only be set for a response if 'data' is as well")
		}
	}
	if j.HTTPStatus < 100 || j.HTTPStatus > 600 {
		return ISE("Response HTTP Status is outside of valid range")
	}

	// probably not the best place for this, but...
	j.JSONAPI.Version = JSONAPIVersion

	return nil
}

// First is just a convenience function that returns the first data object from the
// array
func (j *JSON) First() *Object {
	return j.Data.List[0]
}

// HasData will return true if the JSON document's Data field is set
func (j *JSON) HasData() bool {
	return j.Data != nil && j.Data.List != nil && len(j.Data.List) > 0
}

// HasErrors will return true if the Errors attribute is not nil.
func (j *JSON) HasErrors() bool {
	return j.Errors != nil && len(j.Errors.Objects) > 0
}

// Data is a terrible wrapper around a list that is necessary in order do the
// parsing magic in UnmarshalJSON. If I pass a list literal type in, the data
// disappears after exiting the function.
type Data struct {
	List
}

// UnmarshalJSON allows us to manually decode a list via the json.Unmarshaler
// interface
func (d *Data) UnmarshalJSON(rawData []byte) error {
	log.Printf("raw = %+v\n", string(rawData))

	// replace the list with a parsed version
	isArray, err := dataIsArray(rawData)
	if err != nil {
		return err
	}

	// if we have an array, unmarshal the list in "manually"
	if isArray {
		var list List

		err := json.Unmarshal(rawData, &list)
		if err != nil {
			log.Printf("err = %+v\n", err)
			return err
		}

		d.List = list

	} else {
		// if we have a single object, unmarshal it, and append it to the list
		object := &Object{}

		err = json.Unmarshal(rawData, &object)
		if err != nil {
			return err
		}

		// append our new object to the list
		d.List = List{object}
	}

	return nil
}

// tagMap parses the json document into an abstract type map to be used for
// inferences and type checking
func dataIsArray(raw []byte) (bool, error) {
	var data interface{}

	// first we parse into tags to do some type checking which is necessary to
	// determine whether we are dealing with objects or arrays
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return false, err
	}

	_, isArray := data.([]interface{})
	return isArray, nil
}
