package jshapi

import (
	"fmt"
	"log"
	"os"
	"strings"

	"goji.io"
	"goji.io/pat"
)

// API is used to direct HTTP requests to resources
type API struct {
	*goji.Mux
	prefix    string
	Resources map[string]*Resource
	Logger    *log.Logger
}

// New initializes a new top level API Resource Handler. The most basic implementation
// is:
//
//	api := New("", nil)
//
// But also supports prefixing(/<api_prefix>/<routes>) and custom logging via
// log.Logger https://godoc.org/log#Logger:
//
//	api := New("v1", log.New(os.Stdout, "apiV1: ", log.Ldate|log.Ltime|log.Lshortfile))
//
func New(prefix string, logger *log.Logger) *API {

	// ensure that our top level prefix is "/" prefixed
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	if logger == nil {
		logger = log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Logger:    logger,
	}
}

// Add implements mux support for a given resource which is effectively handled as:
// pat.New("/(prefix/)resource.Plu*)
func (a *API) Add(resource *Resource) {

	// ensure the resource is properly prefixed, and has access to the API logger
	resource.prefix = a.prefix
	resource.Logger = a.Logger

	// track our associated resources, will enable auto-generation docs later
	a.Resources[resource.Type] = resource

	// Add resource wild card to the API mux. Use the resources Matcher() function
	// after an API prefix is applied, as it does the dirty work of building the route
	// automatically for us
	a.Mux.HandleC(pat.New(resource.Matcher()+"*"), resource)
}

// RouteTree prints out all accepted routes for the API that use jshapi implemented
// ways of adding routes through resources: NewCRUDResource(), .Get(), .Post, .Delete(),
// .Patch(), .List(), and .NewAction()
func (a *API) RouteTree() string {
	var routes string

	for _, resource := range a.Resources {
		routes = strings.Join([]string{routes, resource.RouteTree()}, "")
	}

	return routes
}
