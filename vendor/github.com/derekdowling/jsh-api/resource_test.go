package jshapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResource(t *testing.T) {

	Convey("Resource Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "foo"
		resource := NewMockResource("api", resourceType, 2, attrs)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		So(len(resource.Routes), ShouldEqual, 5)

		Convey("->Matcher()", func() {
			resource.prefix = "/api"
			So(resource.Matcher(), ShouldEqual, "/api/"+resourceType+"s")
		})

		Convey("->IDMatcher()", func() {
			resource.prefix = "/api"
			So(resource.IDMatcher(), ShouldEqual, "/api/"+resourceType+"s/:id")
		})

		Convey("->Post()", func() {
			object := sampleObject("", resourceType, attrs)
			object, resp, err := jsc.Post(baseURL, object)

			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			So(object.ID, ShouldEqual, "1")
		})

		Convey("->List()", func() {
			list, resp, err := jsc.GetList(baseURL, resourceType)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)

			So(len(list), ShouldEqual, 2)
			So(list[0].ID, ShouldEqual, "1")
		})

		Convey("->Get()", func() {
			object, resp, err := jsc.GetObject(baseURL, resourceType, "3")

			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(object.ID, ShouldEqual, "3")
		})

		Convey("->Patch()", func() {
			object := sampleObject("1", resourceType, attrs)
			object, resp, err := jsc.Patch(baseURL, object)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(object.ID, ShouldEqual, "1")
		})

		Convey("->Delete()", func() {
			resp, err := jsc.Delete(baseURL, resourceType, "1")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})
	})
}

func TestMutateHandler(t *testing.T) {

	Convey("Custom Handler Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource("/foo", resourceType, 2, attrs)

		handler := func(ctx context.Context, id string) (*jsh.Object, *jsh.Error) {
			object := sampleObject(id, resourceType, attrs)
			return object, nil
		}

		resource.Mutate("mutate", handler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + resource.IDMatcher()

		Convey("Resource State", func() {
			So(len(resource.Routes), ShouldEqual, 6)
			So(resource.Routes[len(resource.Routes)-1], ShouldEqual, "PATCH - /foo/bars/:id/mutate")
		})

		Convey("->Custom()", func() {
			response, err := jsc.Get(baseURL + "/mutate")
			So(err, ShouldBeNil)

			_, err = jsc.ParseObject(response)
			So(err, ShouldBeNil)
		})
	})
}

func TestToOne(t *testing.T) {

	Convey("Relationship ToOne Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource("/foo", resourceType, 2, attrs)

		relationshipHandler := func(ctx context.Context, resourceID string) (*jsh.Object, *jsh.Error) {
			return sampleObject("1", "baz", map[string]string{"baz": "ball"}), nil
		}

		subResourceType := "baz"
		resource.ToOne(subResourceType, relationshipHandler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + resource.IDMatcher()

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/baz", func() {
				resp, err := jsc.Get(baseURL + "/" + subResourceType)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)

				object, err := jsc.ParseObject(resp)
				So(err, ShouldBeNil)
				So(object.ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/baz", func() {
				resp, err := jsc.Get(baseURL + "/relationships/" + subResourceType)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)

				obj, err := jsc.ParseObject(resp)
				So(err, ShouldBeNil)
				So(obj.ID, ShouldEqual, "1")
			})
		})
	})
}

func TestToMany(t *testing.T) {

	Convey("Relationship ToMany Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource("/foo", resourceType, 2, attrs)

		relationshipHandler := func(ctx context.Context, resourceID string) (jsh.List, *jsh.Error) {
			return jsh.List{
				sampleObject("1", "baz", map[string]string{"baz": "ball"}),
				sampleObject("2", "baz", map[string]string{"baz": "ball2"}),
			}, nil
		}

		subResourceType := "baz"
		resource.ToMany(subResourceType, relationshipHandler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + resource.IDMatcher()

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/bazs", func() {
				resp, err := jsc.Get(baseURL + "/" + subResourceType + "s")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)

				list, err := jsc.ParseList(resp)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 2)
				So(list[0].ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/bazs", func() {
				resp, err := jsc.Get(baseURL + "/relationships/" + subResourceType + "s")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)

				list, err := jsc.ParseList(resp)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 2)
				So(list[0].ID, ShouldEqual, "1")
			})
		})

	})
}
