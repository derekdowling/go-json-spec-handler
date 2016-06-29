package jshapi

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestResource(t *testing.T) {

	resource := NewMockResource(testResourceType, 2, testObjAttrs)

	api := New("")
	api.Add(resource)

	server := httptest.NewServer(api)
	baseURL := server.URL

	routeCount := len(resource.Routes)
	if routeCount != 5 {
		log.Fatalf("Invalid number of base resource routes: %d", routeCount)
	}

	Convey("Resource Tests", t, func() {

		Convey("->NewResource()", func() {

			Convey("should be agnostic to plurality", func() {
				resource := NewResource("users")
				So(resource.Type, ShouldEqual, "users")

				resource2 := NewResource("user")
				So(resource2.Type, ShouldEqual, "user")
			})
		})

		Convey("->Post()", func() {
			object := sampleObject("", testResourceType, testObjAttrs)
			doc, resp, err := jsc.Post(baseURL, object)

			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			So(err, ShouldBeNil)
			So(doc.Data[0].ID, ShouldEqual, "1")
		})

		Convey("->List()", func() {
			doc, resp, err := jsc.List(baseURL, testResourceType)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(len(doc.Data), ShouldEqual, 2)
			So(doc.Data[0].ID, ShouldEqual, "1")
		})

		Convey("->Fetch()", func() {
			doc, resp, err := jsc.Fetch(baseURL, testResourceType, "3")

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(doc.Data[0].ID, ShouldEqual, "3")
		})

		Convey("->Patch()", func() {

			Convey("should reject requests with ID mismatch", func() {
				object := sampleObject("1", testResourceType, testObjAttrs)
				request, err := jsc.PatchRequest(baseURL, object)
				So(err, ShouldBeNil)
				// Manually replace resource ID in URL to be invalid
				request.URL.Path = strings.Replace(request.URL.Path, "1", "2", 1)
				doc, resp, err := jsc.Do(request, jsh.ObjectMode)

				So(resp.StatusCode, ShouldEqual, 422)
				So(err, ShouldBeNil)
				So(doc, ShouldNotBeNil)
			})

			Convey("should accept patch requests", func() {
				object := sampleObject("1", testResourceType, testObjAttrs)
				doc, resp, err := jsc.Patch(baseURL, object)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})
		})

		Convey("->Delete()", func() {
			resp, err := jsc.Delete(baseURL, testResourceType, "1")

			So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
			So(err, ShouldBeNil)
		})
	})
}

func TestActionHandler(t *testing.T) {

	resource := NewMockResource(testResourceType, 2, testObjAttrs)

	// Add our custom action
	handler := func(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
		object := sampleObject(id, testResourceType, testObjAttrs)
		return object, nil
	}
	resource.Action("testAction", handler)

	api := New("")
	api.Add(resource)

	server := httptest.NewServer(api)
	baseURL := server.URL

	Convey("Action Handler Tests", t, func() {

		Convey("Resource State", func() {
			So(len(resource.Routes), ShouldEqual, 6)
			So(resource.Routes[len(resource.Routes)-1], ShouldEqual, "PATCH - /bars/:id/testAction")
		})

		Convey("->Custom()", func() {
			doc, response, err := jsc.Action(baseURL, testResourceType, "1", "testAction")

			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, http.StatusOK)
			So(doc.Data, ShouldNotBeEmpty)
		})
	})
}

func TestToOne(t *testing.T) {

	resource := NewMockResource(testResourceType, 2, testObjAttrs)

	relationshipHandler := func(ctx context.Context, resourceID string) (*jsh.Object, jsh.ErrorType) {
		return sampleObject("1", "baz", map[string]string{"baz": "ball"}), nil
	}

	subResourceType := "baz"
	resource.ToOne(subResourceType, relationshipHandler)

	api := New("")
	api.Add(resource)

	server := httptest.NewServer(api)
	baseURL := server.URL

	Convey("Relationship ToOne Tests", t, func() {

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/baz", func() {
				doc, resp, err := jsc.Action(baseURL, testResourceType, "1", subResourceType)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/baz", func() {
				doc, resp, err := jsc.Action(baseURL, testResourceType, "1", "relationships/"+subResourceType)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})
		})
	})
}

func TestToMany(t *testing.T) {

	resource := NewMockResource(testResourceType, 2, testObjAttrs)

	relationshipHandler := func(ctx context.Context, resourceID string) (jsh.List, jsh.ErrorType) {
		return jsh.List{
			sampleObject("1", "baz", map[string]string{"baz": "ball"}),
			sampleObject("2", "baz", map[string]string{"baz": "ball2"}),
		}, nil
	}

	subResourceType := "baz"
	resource.ToMany(subResourceType, relationshipHandler)

	api := New("")
	api.Add(resource)

	server := httptest.NewServer(api)
	baseURL := server.URL

	Convey("Relationship ToMany Tests", t, func() {

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/bazs", func() {
				doc, resp, err := jsc.Action(baseURL, testResourceType, "1", subResourceType+"s")

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(len(doc.Data), ShouldEqual, 2)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/bazs", func() {
				doc, resp, err := jsc.Action(baseURL, testResourceType, "1", "relationships/"+subResourceType+"s")

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(len(doc.Data), ShouldEqual, 2)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})
		})
	})
}
