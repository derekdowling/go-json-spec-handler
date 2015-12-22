Go JSON API Specification Handler
---

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/derekdowling/go-json-spec-handler)
[![Travis CI](https://img.shields.io/travis/derekdowling/go-json-spec-handler/master.svg?style=flat-square)](https://travis-ci.org/derekdowling/go-json-spec-handler)
[![Go Report Card](http://goreportcard.com/badge/manyminds/api2go)](http://goreportcard.com/report/derekdowling/go-json-spec-handler)
[TestCoverage](http://gocover.io/github.com/derekdowling/go-json-spec-handler?version=1.5rc1)

A server (de)serialization handler for creating [JSON API Specification](http://jsonapi.org/) 
compatible backends in Go. Works with [Ember-Data](https://github.com/emberjs/data) too!

# Packages

### jsh - JSON Specification Handler

Streamlined JSON input/output handling for a new, or existing API server:

```go
import github.com/derekdowling/go-json-spec-handler

type User struct {
  // valid from github.com/asaskevich/govalidator gives us input validation
  // when object.Unmarshal() is invoked on this type
  Name string `json:"name" valid:"alphanum"`
}

func PatchUser(w http.ResponseWriter, r *http.Request) {
  user := &User{}

  // performs Specification checks against the request
  object, err := jsh.ParseObject(*http.Request)
  if err != nil {
    jsh.Send(w, r, err)
    return
  }

  // use object.ID to look up user/do business logic

  // unmarshal data into relevant internal types if govalidator passes, otherwise
  // return the pre-formatted HTTP 422 error to signify how the input failed
  err := object.Unmarshal("user", user)
  if err != nil {
    jsh.Send(w, r, err)
    return
  }

  // modify, re-package, and send the object
  user.Name = "Bob"
  err := object.Marshal(user)
  if err != nil {
    jsh.Send(w, r, err)
  }

  jsh.Send(w, r, object)
}
```

### jsc - JSON Specification Client

HTTP JSON Client for interacting with JSON APIs. Built on top of http.Client
and jsh.

```go
import github.com/derekdowling/go-json-spec-handler/client

// GET http://your.api/users/1
object, response, err := jsc.GetObject("http://your.api/", "user", "1")
```


### Philosophy Behind JSH

In sticking with Go's philosophy of modules over frameworks, `jsh` was created
to be a drop in serialization layer focusing only on parsing, validating, and
sending JSON API compatible responses. Currently `jsh` is getting fairly close
to stable. It's undergone a number of large refactors to accomodate new
aspects of the specification as I round out the expected feature set which is
pretty well completed, including support for the HTTP client linked above.

If you're looking for a good place to start with a new API, I've since created
[jshapi](https://github.com/derekdowling/jsh-api) which builds on top of [Goji 2](https://goji.io/)
and `jsh` in order to handle the routing structure that JSON API requires as
well as a number of other useful tools for testing and mocking APIs as you
develop your own projects.

### Features 

    Implemented:

    - Handles both single object and array based JSON requests and responses
    - Input validation with HTTP 422 Status support via [go-validator](https://github.com/go-validator/validator)
    - Client request validation with HTTP 406 Status responses
    - Links, Relationship, Meta fields
    - Prepackaged error responses, easy to use Internal Service Error builder
    - Smart responses with correct HTTP Statuses based on Request Method and HTTP Headers
    - HTTP Client for GET, POST, DELETE, PATCH

    TODO:

    - [Reserved character checking](http://jsonapi.org/format/upcoming/#document-member-names-reserved-characters)

    Not Implementing:

    * These features aren't handled because they are beyond the scope of what
      this module is meant to be. See [jshapi](https://github.com/derekdowling/jsh-api)
      if these are problems that you'd also like to have solved.

    - Routing
    - Relationship management
    - Sorting
    - Pagination
    - Filtering

## Examples

There are lots of great examples in the tests themselves that show exactly how jsh works.
The [godocs](https://godoc.org/github.com/derekdowling/go-json-spec-handler) as linked above
have a number of examples in them as well.
