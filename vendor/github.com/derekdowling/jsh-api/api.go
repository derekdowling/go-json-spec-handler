package jshapi

import (
	"log"
	"os"

	"github.com/zenazn/goji/web"
)

// API is used to direct HTTP requests to resources
type API struct {
	*web.Mux
	prefix    string
	Resources map[string]*Resource
	Logger    *log.Logger
}

// New initializes a Handler object
func New(prefix string) *API {
	return &API{
		Mux:       web.New(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Logger:    log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// AddResource adds a new resource of type "name" to the API's router
func (a *API) AddResource(resource *Resource) {

	// add prefix and logger
	resource.prefix = a.prefix
	resource.Logger = a.Logger

	a.Resources[resource.name] = resource

	// Add subrouter to main API mux, use Matcher plus catch all
	a.Mux.Handle(resource.Matcher()+"*", resource.Mux)
}
