Go JSON API Specification Handler
---

[![GoDoc](https://godoc.org/github.com/derekdowling/go-json-spec-handler?status.png)](https://godoc.org/github.com/derekdowling/go-json-spec-handler)
[![Go Report Card](http://goreportcard.com/badge/manyminds/api2go)](http://goreportcard.com/report/derekdowling/go-json-spec-handler)
[TestCoverage](http://gocover.io/github.com/derekdowling/go-json-spec-handler?version=1.5rc1)

An HTTP Client and Server request/response handler for dealing with [JSON Specification](http://jsonapi.org/) 
APIs. Great for Ember.js!

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

    - Reserved character checking

    Not Implenting:

    * These features aren't handled because they are beyond the scope of what
      this library is meant to be. In the future, I might build a framework
      utilizing this library to handle these complex features which require
      Router and ORM compatibility.

    - Relationship management
    - Sorting
    - Pagination
    - Filtering

## Examples

- [jshapi](https://github.com/derekdowling/jsh-api) abstracts the full
  serialization layer for JSON Spec APIs.

There are lots of great examples in the tests themselves that show exactly how jsh works.
The [godocs](https://godoc.org/github.com/derekdowling/go-json-spec-handler) as linked above
have a number of examples in them as well.
