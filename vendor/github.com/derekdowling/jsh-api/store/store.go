// Package store is a collection of composable interfaces that are can be implemented
// in order to build a storage driver
package store

import (
	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
)

// CRUD implements all sub-storage functions
type CRUD interface {
	Save(ctx context.Context, object *jsh.Object) (*jsh.Object, *jsh.Error)
	Get(ctx context.Context, id string) (*jsh.Object, *jsh.Error)
	List(ctx context.Context) (jsh.List, *jsh.Error)
	Update(ctx context.Context, object *jsh.Object) (*jsh.Object, *jsh.Error)
	Delete(ctx context.Context, id string) *jsh.Error
}

// Save a new resource to storage
type Save func(ctx context.Context, object *jsh.Object) (*jsh.Object, *jsh.Error)

// Get a specific instance of a resource by id from storage
type Get func(ctx context.Context, id string) (*jsh.Object, *jsh.Error)

// List all instances of a resource from storage
type List func(ctx context.Context) (jsh.List, *jsh.Error)

// Update an existing object in storage
type Update func(ctx context.Context, object *jsh.Object) (*jsh.Object, *jsh.Error)

// Delete an object from storage by id
type Delete func(ctx context.Context, id string) *jsh.Error

// ToMany retrieves a list of objects of a single resource type that are related to
// the provided resource id
type ToMany func(ctx context.Context, id string) (jsh.List, *jsh.Error)
