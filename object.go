package jsh

import (
	"encoding/json"
	"fmt"
)

// Object represents the default JSON spec for objects
type Object struct {
	Type       string          `json:"type"`
	ID         string          `json:"id"`
	Attributes json.RawMessage `json:"attributes"`
}

// NewObject prepares a new JSON Object for an API response. Whatever is provided
// as attributes will be marshalled to JSON.
func NewObject(id string, objType string, attributes interface{}) (*Object, error) {
	object := &Object{
		ID:   id,
		Type: objType,
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
func (o *Object) Unmarshal(objType string, target interface{}) error {

	if objType != o.Type {
		return fmt.Errorf("Attempting to convert object to incompatible type")
	}

	err := json.Unmarshal(o.Attributes, target)
	if err != nil {
		return fmt.Errorf(
			"Error converting %s to type '%s': %s",
			o.Attributes,
			objType,
			err.Error(),
		)
	}

	return nil
}
