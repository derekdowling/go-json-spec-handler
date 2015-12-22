package jshapi

import (
	"log"

	"github.com/derekdowling/go-json-spec-handler"
)

const testType = "test"

// NewMockResource builds a mock API endpoint that can perform basic CRUD actions:
//
//	GET    /types
//	POST   /types
//	GET    /types/:id
//	DELETE /types/:id
//	PATCH  /types/:id
//
// Will return objects and lists based upon the sampleObject that is specified here
// in the constructor.
func NewMockResource(resourceType string, listCount int, sampleObject interface{}) *Resource {
	mock := &MockStorage{
		ResourceType:       resourceType,
		ResourceAttributes: sampleObject,
		ListCount:          listCount,
	}

	return NewCRUDResource(resourceType, mock)
}

func sampleObject(id string, resourceType string, sampleObject interface{}) *jsh.Object {
	object, err := jsh.NewObject(id, resourceType, sampleObject)
	if err != nil {
		log.Fatal(err.Error())
	}

	return object
}
