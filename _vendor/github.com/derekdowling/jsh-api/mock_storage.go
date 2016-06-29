package jshapi

import (
	"log"
	"strconv"

	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
)

// MockStorage allows you to mock out APIs really easily, and is also used internally
// for testing the API layer.
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
func (m *MockStorage) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	var err *jsh.Error
	object.ID = "1"

	return object, err
}

// Get returns a resource with ID as specified by the request
func (m *MockStorage) Get(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
	var err *jsh.Error

	return m.SampleObject(id), err
}

// List returns a sample list
func (m *MockStorage) List(ctx context.Context) (jsh.List, jsh.ErrorType) {
	var err *jsh.Error

	return m.SampleList(m.ListCount), err
}

// Update does nothing
func (m *MockStorage) Update(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	var err jsh.ErrorList
	err = nil

	return object, err
}

// Delete does nothing
func (m *MockStorage) Delete(ctx context.Context, id string) jsh.ErrorType {
	var err *jsh.Error

	return err
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
