package jshapi

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/derekdowling/goji2-logger"

	"goji.io"
	"goji.io/pat"
)

// Logger is the default logging interface used in JSH API
type Logger gojilogger.Logger

// API is used to direct HTTP requests to resources
type API struct {
	*goji.Mux
	prefix    string
	Resources map[string]*Resource
	Logger    Logger
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
func New(prefix string, logger Logger) *API {

	// ensure that our top level prefix is "/" prefixed
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	if logger == nil {
		logger = log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// create our new logger
	api := &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Logger:    logger,
	}

	// register default middleware
	gojilogger.SetLogger(logger)
	api.UseC(gojilogger.Middleware)

	return api
}

// Add implements mux support for a given resource which is effectively handled as:
// pat.New("/(prefix/)resource.Plu*)
func (a *API) Add(resource *Resource) {

	// ensure the resource is properly prefixed, and has access to the API logger
	resource.Logger = a.Logger

	// track our associated resources, will enable auto-generation docs later
	a.Resources[resource.Type] = resource

	// Because of how prefix matches work:
	// https://godoc.org/github.com/goji/goji/pat#hdr-Prefix_Matches
	// We need two separate routes,
	// /(prefix/)resources
	matcher := path.Join(a.prefix, resource.PluralType())
	a.Mux.HandleC(pat.New(matcher), resource)

	// And:
	// /(prefix/)resources/*
	idMatcher := path.Join(a.prefix, resource.PluralType(), "*")
	a.Mux.HandleC(pat.New(idMatcher), resource)
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
