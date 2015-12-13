# JSH-API

[![GoDoc](https://godoc.org/github.com/derekdowling/go-json-spec-handler?status.png)](https://godoc.org/github.com/derekdowling/jsh-api) [![Build Status](https://travis-ci.org/derekdowling/jsh-api.svg?branch=master)](https://travis-ci.org/derekdowling/jsh-api)

A [JSON Specification](http://jsonapi.org) API Build created on top of
[jsh](http://github.com/derekdowling/go-json-spec-handler). Bring your own
router, bring your own storage, focus on functionality, and let jsh-api do the
rest.

## Setup

```go
import github.com/derekdowling/jsh-api

api := jshapi.New("<prefix>")

// implement jshapi/store.CRUD interface, then:
userStorage := &UserStorage{}
resource := jshapi.NewCRUDResource("user", userStorage)
resource.Use(yourUserMiddleware)

// add resources to the API
api.AddResource(resource)

// API middleware
api.Use(yourTopLevelAPIMiddleware)
http.ListenAndServe("localhost:8000", api)
```

## What It Handles

All of the dirty work for parsing all supported JSON API request endpoints for
each resource:

```
POST /resources
GET /resources
GET /resources/:id
DELETE /resources/:id
PATCH /resources/:id
```

## Implementing a Storage Driver with jsh

Below is a simple example of how you might implement the required 
[Storage Driver](https://godoc.org/github.com/derekdowling/jsh-api/store#CRUD) for a
given resource using
[jsh](https://godoc.org/github.com/derekdowling/go-json-spec-handler) for Save
and Update. This should give you a pretty good idea of how easy it is to
implement the Storage driver with jsh.


```go
type User struct {
    ID string
    Name string `json:"name"`
}

func Save(ctx context.Context, object *jsh.Object) jsh.SendableError {
    user := &User{}
    err := object.Unmarshal("user", user)
    if err != nil {
        return err
    }

    // generate your id, however you choose
    user.ID = "1234"

    // do save logic
    return nil
}

func Update(ctx context.Context, object *jsh.Object) jsh.SendableError {
    user := &User{}
    err := object.Unmarshal("user", user)
    if err != nil {
        return err
    }

    // object has the lookup ID
    id := object.ID

    // perform patch
    return nil
}
```
