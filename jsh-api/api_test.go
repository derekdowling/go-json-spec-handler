package jshapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

const testResourceType = "bars"

func TestAPI(t *testing.T) {

	Convey("API Tests", t, func() {

		api := New("api")

		So(api.prefix, ShouldEqual, "/api")

		testAttrs := map[string]string{
			"foo": "bar",
		}

		Convey("->AddResource()", func() {
			resource := NewMockResource(testResourceType, 1, testAttrs)
			api.Add(resource)

			So(api.Resources[testResourceType], ShouldEqual, resource)

			server := httptest.NewServer(api)
			baseURL := server.URL + api.prefix

			Convey("should work with /<resource> routes", func() {
				_, resp, err := jsc.List(baseURL, testResourceType)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
			})

			Convey("should work with /<resource>/:id routes", func() {
				patchObj, err := jsh.NewObject("1", testResourceType, testAttrs)
				So(err, ShouldBeNil)

				_, resp, patchErr := jsc.Patch(baseURL, patchObj)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(patchErr, ShouldBeNil)
			})
		})
	})
}
