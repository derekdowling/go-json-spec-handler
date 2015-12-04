package jshapi

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/zenazn/goji/web"
)

// Resource is a handler object for dealing with CRUD API endpoints
type Resource struct {
	*web.Mux
	// The singular name of the resource type ex) "user" or "post"
	name string
	// An implemented jshapi.Storage interface
	storage Storage
	// An implementation of Go's standard logger
	Logger *log.Logger
	// Prefix is set if the resource is not the top level of URI, "/prefix/resources
	prefix string
}

// NewResource is a resource constructor
func NewResource(prefix string, name string, storage Storage) *Resource {

	r := &Resource{
		Mux:     web.New(),
		name:    name,
		storage: storage,
		prefix:  prefix,
	}

	// setup resource sub-router
	r.Mux.Post(r.Matcher(), r.Post)
	r.Mux.Get(r.IDMatcher(), r.Get)
	r.Mux.Get(r.Matcher(), r.List)
	r.Mux.Delete(r.IDMatcher(), r.Delete)
	r.Mux.Patch(r.IDMatcher(), r.Patch)

	return r
}

// Post => POST /resources
func (res *Resource) Post(c web.C, w http.ResponseWriter, r *http.Request) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	err = res.storage.Save(object)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

// Get => GET /resources/:id
func (res *Resource) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Printf("r.URL = %+v\n", r.URL)
	id, exists := c.URLParams["id"]
	if !exists {
		res.sendAndLog(c, w, r, jsh.ISE(fmt.Sprintf("Unable to parse resource ID from path: %s", r.URL.Path)))
		return
	}

	object, err := res.storage.Get(id)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

// List => GET /resources
func (res *Resource) List(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Printf("r.URL = %+v\n", r.URL)
	list, err := res.storage.List()
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, list)
}

// Delete => DELETE /resources/:id
func (res *Resource) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	id, exists := c.URLParams["id"]
	if !exists {
		res.sendAndLog(c, w, r, jsh.ISE(fmt.Sprintf("Unable to parse resource ID from path: %s", r.URL.Path)))
		return
	}

	err := res.storage.Delete(id)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Patch => PATCH /resources/:id
func (res *Resource) Patch(c web.C, w http.ResponseWriter, r *http.Request) {
	object, err := jsh.ParseObject(r)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	err = res.storage.Patch(object)
	if err != nil {
		res.sendAndLog(c, w, r, err)
		return
	}

	res.sendAndLog(c, w, r, object)
}

func (res *Resource) sendAndLog(c web.C, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {
	jshErr, isType := sendable.(*jsh.Error)
	if isType && jshErr.Status == http.StatusInternalServerError {
		res.Logger.Printf("JSH ISE: %s-%s", jshErr.Title, jshErr.Detail)
	}

	err := jsh.Send(w, r, sendable)
	if err != nil {
		res.Logger.Print(err.Error())
	}
}

// PluralType returns the resource's name, but pluralized
func (res *Resource) PluralType() string {
	return res.name + "s"
}

// IDMatcher returns a uri path matcher for the resource type
func (res *Resource) IDMatcher() string {
	return path.Join(res.Matcher(), ":id")
}

// Matcher returns the top level uri path matcher for the resource type
func (res *Resource) Matcher() string {
	return fmt.Sprintf("/%s", path.Join(res.prefix, res.PluralType()))
}
