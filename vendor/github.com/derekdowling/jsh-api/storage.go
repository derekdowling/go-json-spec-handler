package jshapi

import "github.com/derekdowling/go-json-spec-handler"

// Storage is an interface that allows jshapi to perform CRUD actions
type Storage interface {
	// Save a new resource to storage
	Save(object *jsh.Object) jsh.SendableError
	// Get a specific instance of a resource from storage
	Get(id string) (*jsh.Object, jsh.SendableError)
	// List all instances of a resource from storage
	List() (jsh.List, jsh.SendableError)
	// Save an object to storage
	Patch(object *jsh.Object) jsh.SendableError
	// Delete from storage
	Delete(id string) jsh.SendableError
}
