package jshapi

import (
	"log"
	"strconv"

	"github.com/derekdowling/go-json-spec-handler"
)

const testType = "test"

// MockStorage allows you to mock out APIs really easily
type MockStorage struct {
	// ResourceType is the name of the resource you are mocking i.e. "user", "comment"
	ResourceType string
	// ResourceAttributes a sample set of attributes a resource object should have
	// used by GET /resources and GET /resources/:id
	ResourceAttributes interface{}
	// ListCount is the number of sample objects to return in a GET /resources request
	ListCount int
}

// Save assigns a URL of 1 to the object
func (m *MockStorage) Save(object *jsh.Object) jsh.SendableError {
	object.ID = "1"
	return nil
}

// Get returns a resource with ID as specified by the request
func (m *MockStorage) Get(id string) (*jsh.Object, jsh.SendableError) {
	return m.SampleObject(id), nil
}

// List returns a sample list
func (m *MockStorage) List() (jsh.List, jsh.SendableError) {
	return m.SampleList(m.ListCount), nil
}

// Patch does nothing
func (m *MockStorage) Patch(object *jsh.Object) jsh.SendableError {
	return nil
}

// Delete does nothing
func (m *MockStorage) Delete(id string) jsh.SendableError {
	return nil
}

// SampleObject builds an object based on provided resource specifications
func (m *MockStorage) SampleObject(id string) *jsh.Object {
	object, err := jsh.NewObject(id, m.ResourceType, m.ResourceAttributes)
	if err != nil {
		log.Fatal(err.Error())
	}

	return object
}

// SampleList generates a sample list of resources that can be used for/against the
// mock API
func (m *MockStorage) SampleList(length int) jsh.List {

	list := jsh.List{}

	for id := 1; id <= length; id++ {
		list = append(list, m.SampleObject(strconv.Itoa(id)))
	}

	return list
}

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
func NewMockResource(prefix string, resourceType string, listCount int, sampleObject interface{}) *Resource {

	mock := &MockStorage{
		ResourceType:       resourceType,
		ResourceAttributes: sampleObject,
		ListCount:          listCount,
	}

	return NewResource("", resourceType, mock)
}

func sampleObject(id string, resourceType string, sampleObject interface{}) *jsh.Object {
	object, err := jsh.NewObject(id, resourceType, sampleObject)
	if err != nil {
		log.Fatal(err.Error())
	}

	return object
}
