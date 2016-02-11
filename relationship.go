package jsh

import (
	"fmt"

	"encoding/json"
)

// Relationship represents a reference from the resource object in which it's
// defined to other resource objects.
type Relationship struct {
	Links *Links                 `json:"links,omitempty"`
	Data  ResourceLinkage        `json:"data,omitempty"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// ResourceLinkage is a typedef around a slice of resource identifiers. This
// allows us to implement a custom UnmarshalJSON.
type ResourceLinkage []*ResourceIdentifier

// ResourceIdentifier identifies an individual resource.
type ResourceIdentifier struct {
	Type string `json:"type" valid:"required"`
	ID   string `json:"id" valid:"required"`
}

/*
UnmarshalJSON allows us to manually decode a the resource linkage via the
json.Unmarshaler interface.
*/
func (rl *ResourceLinkage) UnmarshalJSON(data []byte) error {
	// Create a sub-type here so when we call Unmarshal below, we don't recursively
	// call this function over and over
	type UnmarshalLinkage ResourceLinkage

	// if our "List" is a single object, modify the JSON to make it into a list
	// by wrapping with "[ ]"
	if data[0] == '{' {
		data = []byte(fmt.Sprintf("[%s]", data))
	}

	newLinkage := UnmarshalLinkage{}

	err := json.Unmarshal(data, &newLinkage)
	if err != nil {
		return err
	}

	convertedLinkage := ResourceLinkage(newLinkage)
	*rl = convertedLinkage

	return nil
}
