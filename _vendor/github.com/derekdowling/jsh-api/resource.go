package jshapi

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"

	"goji.io"
	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api/store"
)

const (
	post    = "POST"
	get     = "GET"
	list    = "LIST"
	delete  = "DELETE"
	patch   = "PATCH"
	patID   = "/:id"
	patRoot = ""
)

/*
Resource holds the necessary state for creating a REST API endpoint for a
given resource type. Will be accessible via `/<type>`.

Using NewCRUDResource you can generate a generic CRUD handler for a
JSON Specification Resource end point. If you wish to only implement a subset
of these endpoints that is also available through NewResource() and manually
registering storage handlers via .Post(), .Get(), .List(), .Patch(), and .Delete():

Besides the built in registration helpers, it isn't recommended, but you can add
your own routes using the goji.Mux API:

	func searchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		name := pat.Param(ctx, "name")
		fmt.Fprintf(w, "Hello, %s!", name)
	}

	resource := jshapi.NewCRUDResource("users", userStorage)
	// creates /users/search/:name
	resource.HandleC(pat.New("search/:name"), searchHandler)
*/
type Resource struct {
	*goji.Mux
	// The singular name of the resource type("user", "post", etc)
	Type string
	// Routes is a list of routes registered to the resource
	Routes []string
	// Map of relationships
	Relationships map[string]Relationship
}

/*
NewResource is a resource constructor that makes no assumptions about routes
that you'd like to implement, but still provides some basic utilities for
managing routes and handling API calls.

The prefix parameter causes all routes created within the resource to be prefixed.
*/
func NewResource(resourceType string) *Resource {
	return &Resource{
		// Mux is a goji.SubMux, inherits context from parent Mux
		Mux: goji.SubMux(),
		// Type of the resource, makes no assumptions about plurality
		Type:          resourceType,
		Relationships: map[string]Relationship{},
		// A list of registered routes, useful for debugging
		Routes: []string{},
	}
}

// NewCRUDResource generates a resource
func NewCRUDResource(resourceType string, storage store.CRUD) *Resource {
	resource := NewResource(resourceType)
	resource.CRUD(storage)
	return resource
}

/*
CRUD is syntactic sugar and a shortcut for registering all JSON API CRUD
routes for a compatible storage implementation:

Registers handlers for:
	GET    /resource
	POST   /resource
	GET    /resource/:id
	DELETE /resource/:id
	PATCH  /resource/:id
*/
func (res *Resource) CRUD(storage store.CRUD) {
	res.Get(storage.Get)
	res.Patch(storage.Update)
	res.Post(storage.Save)
	res.List(storage.List)
	res.Delete(storage.Delete)
}

// Post registers a `POST /resource` handler with the resource
func (res *Resource) Post(storage store.Save) {
	res.HandleFuncC(
		pat.Post(patRoot),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.postHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(post, patRoot)
}

// Get registers a `GET /resource/:id` handler for the resource
func (res *Resource) Get(storage store.Get) {
	res.HandleFuncC(
		pat.Get(patID),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.getHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(get, patID)
}

// List registers a `GET /resource` handler for the resource
func (res *Resource) List(storage store.List) {
	res.HandleFuncC(
		pat.Get(patRoot),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.listHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(get, patRoot)
}

// Delete registers a `DELETE /resource/:id` handler for the resource
func (res *Resource) Delete(storage store.Delete) {
	res.HandleFuncC(
		pat.Delete(patID),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.deleteHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(delete, patID)
}

// Patch registers a `PATCH /resource/:id` handler for the resource
func (res *Resource) Patch(storage store.Update) {
	res.HandleFuncC(
		pat.Patch(patID),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.patchHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(patch, patID)
}

// ToOne registers a `GET /resource/:id/(relationships/)<resourceType>` route which
// returns a "resourceType" in a One-To-One relationship between the parent resource
// type and "resourceType" as specified here. The "/relationships/" uri component is
// optional.
//
// CRUD actions on a specific relationship "resourceType" object should be performed
// via it's own top level /<resourceType> jsh-api handler as per JSONAPI specification.
func (res *Resource) ToOne(
	resourceType string,
	storage store.Get,
) {
	resourceType = strings.TrimSuffix(resourceType, "s")

	res.relationshipHandler(
		resourceType,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.getHandler(ctx, w, r, storage)
		},
	)

	res.Relationships[resourceType] = ToOne
}

// ToMany registers a `GET /resource/:id/(relationships/)<resourceType>s` route which
// returns a list of "resourceType"s in a One-To-Many relationship with the parent resource.
// The "/relationships/" uri component is optional.
//
// CRUD actions on a specific relationship "resourceType" object should be performed
// via it's own top level /<resourceType> jsh-api handler as per JSONAPI specification.
func (res *Resource) ToMany(
	resourceType string,
	storage store.ToMany,
) {
	if !strings.HasSuffix(resourceType, "s") {
		resourceType = fmt.Sprintf("%ss", resourceType)
	}

	res.relationshipHandler(
		resourceType,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.toManyHandler(ctx, w, r, storage)
		},
	)

	res.Relationships[resourceType] = ToMany
}

// relationshipHandler does the dirty work of setting up both routes for a single
// relationship
func (res *Resource) relationshipHandler(
	resourceType string,
	handler goji.HandlerFunc,
) {

	// handle /.../:id/<resourceType>
	matcher := fmt.Sprintf("%s/%s", patID, resourceType)
	res.HandleFuncC(
		pat.Get(matcher),
		handler,
	)
	res.addRoute(get, matcher)

	// handle /.../:id/relationships/<resourceType>
	relationshipMatcher := fmt.Sprintf("%s/relationships/%s", patID, resourceType)
	res.HandleFuncC(
		pat.Get(relationshipMatcher),
		handler,
	)
	res.addRoute(get, relationshipMatcher)
}

// Action allows you to add custom actions to your resource types, it uses the
// GET /(prefix/)resourceTypes/:id/<actionName> path format
func (res *Resource) Action(actionName string, storage store.Get) {
	matcher := path.Join(patID, actionName)

	res.HandleFuncC(
		pat.Get(matcher),
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			res.actionHandler(ctx, w, r, storage)
		},
	)

	res.addRoute(patch, matcher)
}

// POST /resources
func (res *Resource) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Save) {
	parsedObject, parseErr := jsh.ParseObject(r)
	if parseErr != nil && reflect.ValueOf(parseErr).IsNil() == false {
		SendHandler(ctx, w, r, parseErr)
		return
	}

	object, err := storage(ctx, parsedObject)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, object)
}

// GET /resources/:id
func (res *Resource) getHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Get) {
	id := pat.Param(ctx, "id")

	object, err := storage(ctx, id)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, object)
}

// GET /resources
func (res *Resource) listHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.List) {
	list, err := storage(ctx)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, list)
}

// DELETE /resources/:id
func (res *Resource) deleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Delete) {
	id := pat.Param(ctx, "id")

	err := storage(ctx, id)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /resources/:id
func (res *Resource) patchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Update) {
	parsedObject, parseErr := jsh.ParseObject(r)
	if parseErr != nil && reflect.ValueOf(parseErr).IsNil() == false {
		SendHandler(ctx, w, r, parseErr)
		return
	}

	object, err := storage(ctx, parsedObject)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, object)
}

// GET /resources/:id/(relationships/)<resourceType>s
func (res *Resource) toManyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.ToMany) {
	id := pat.Param(ctx, "id")

	list, err := storage(ctx, id)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, list)
}

// All HTTP Methods for /resources/:id/<mutate>
func (res *Resource) actionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, storage store.Get) {
	id := pat.Param(ctx, "id")

	response, err := storage(ctx, id)
	if err != nil && reflect.ValueOf(err).IsNil() == false {
		SendHandler(ctx, w, r, err)
		return
	}

	SendHandler(ctx, w, r, response)
}

// addRoute adds the new method and route to a route Tree for debugging and
// informational purposes.
func (res *Resource) addRoute(method string, route string) {
	res.Routes = append(res.Routes, fmt.Sprintf("%s - /%s%s", method, res.Type, route))
}

// RouteTree prints a recursive route tree based on what the resource, and
// all subresources have registered
func (res *Resource) RouteTree() string {
	var routes string

	for _, route := range res.Routes {
		routes = strings.Join([]string{routes, route}, "\n")
	}

	return routes
}
