Go JSON API Specification Handler
---

[![GoDoc](https://godoc.org/github.com/derekdowling/go-json-spec-handler?status.png)](https://godoc.org/github.com/derekdowling/go-json-spec-handler) [![Build Status](https://img.shields.io/travis/derekdowling/go-json-spec-handler.svg)](https://travis-ci.org/derekdowling/go-json-spec-handler)

A Golang API helper that deals with request serialization and response sending for creating a [JSON API Specification](http://jsonapi.org/) compatible Golang API. Great for Ember.js!

## Features

    Implemented:

    - Handles both single object and array based JSON requests and responses
    - Input validation with HTTP 422 Status support via [go-validator](https://github.com/go-validator/validator)
    - Client request validation with HTTP 406 Status responses
    - Links, Relationship, Meta fields
    - Prepackaged error responses, easy to use Internal Service Error builder
    - Smart responses with correct HTTP Statuses based on Request Method and HTTP Headers

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

## Installation

```
$ go get github.com/derekdowling/go-json-spec-handler
```

## Examples

There are lots of great examples in the tests themselves that show exactly how it works, also check out the [godocs](https://godoc.org/github.com/derekdowling/go-json-spec-handler) as linked above.
