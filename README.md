A Go JSONAPI Specification Handler
---

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/derekdowling/go-json-spec-handler)
[![Travis CI](https://img.shields.io/travis/derekdowling/go-json-spec-handler/master.svg?style=flat-square)](https://travis-ci.org/derekdowling/go-json-spec-handler)
[![Go Report Card](http://goreportcard.com/badge/manyminds/api2go)](http://goreportcard.com/report/derekdowling/go-json-spec-handler)
[TestCoverage](http://gocover.io/github.com/derekdowling/go-json-spec-handler?version=1.5rc1)

A (de)serialization handler for writing [JSON API Specification](http://jsonapi.org/) 
compatible software in Go. Works with [Ember/Ember-Data](https://github.com/emberjs/data) too!

# Contents

1. [JSH](#jsh---json-specification-handler)
  * [Motivation](#motivation-for-jsh)
  * [Features](#features)
  * [Stability](#stability)
2. [JSC](#jsc---json-specification-client)
3. [JSH-API](#jsh-api)

### jsh - JSON Specification Handler

For streamlined JSONAPI object serialization. Uses [govalidator](github.com/asaskevich/govalidator) for input validation.

```go
import github.com/derekdowling/go-json-spec-handler

type User struct {
  // valid from github.com/asaskevich/govalidator gives us input validation
  // when object.Unmarshal() is invoked on this type
  Name string `json:"name" valid:"alphanum"`
}

// example http.HandlerFunc
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
  err = object.Unmarshal("user", user)
  if err != nil {
    jsh.Send(w, r, err)
    return
  }

  // modify your internal type
  user.Name = "Bob"

  // repackage and send the JSONAPI object
  err = object.Marshal(user)
  if err != nil {
    jsh.Send(w, r, err)
  }

  jsh.Send(w, r, object)
}
```

### Motivation for JSH

JSH was written for tackling the issue of dealing with Ember-Data within a pre-existing
API server. In sticking with Go's philosophy of modules over frameworks, it is intended
to be a drop in serialization layer focusing only on parsing, validating, and
sending JSONAPI compatible responses.

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
      this module is meant to achieve. See [jshapi](https://github.com/derekdowling/jsh-api)
      for a full-fledged API solution that solves many of these problems.

    - Routing
    - Sorting
    - Pagination
    - Filtering
    - ORM

### Stability

`jsh` has a mostly stabilized core data document model. At this point in time I am not yet
ready to declare v1, but am actively trying to avoid breaking the public API. The areas most
likely to receive improvement include relationship, link, and metadata management. At this
point in time I can confidentally suggest you use `jsh` without risking major upgrade incompatibility
going forward!


### [jsc - JSON Specification Client](https://godoc.org/github.com/derekdowling/go-json-spec-handler/client)

A HTTP JSONAPI Client for making outbound server requests. Built on top of http.Client and jsh:

```go
import github.com/derekdowling/go-json-spec-handler/client

// GET http://your.api/users/1
document, resp, err := jsc.Fetch("http://your.api/", "users", "1")
object := doc.First()

user := &yourUser{}
err := object.Unmarshal("users", user)
```

### [JSH-API](https://github.com/derekdowling/jsh-api)

If you're looking for a good place to start with a new API, I've since created
[jshapi](https://github.com/derekdowling/jsh-api) which builds on top of [Goji 2](https://goji.io/)
and `jsh` in order to handle the routing structure that JSON API requires as
well as a number of other useful tools for testing and mocking APIs as you
develop your own projects.

## Examples

There are lots of great examples in the tests themselves that show exactly how jsh works.
The [godocs](https://godoc.org/github.com/derekdowling/go-json-spec-handler) as linked above
have a number of examples in them as well.
