package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/jsh-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	Convey("Get Tests", t, func() {

		resource := jshapi.NewMockResource("test", 1, nil)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Get()", func() {
			json, resp, err := Get(baseURL + "/tests/1")

			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(json.HasErrors(), ShouldBeFalse)
			So(json.HasData(), ShouldBeTrue)
		})

		Convey("->GetList()", func() {

			Convey("should handle an object listing request", func() {
				json, resp, err := List(baseURL, "test")

				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(json.HasErrors(), ShouldBeFalse)
				So(json.HasData(), ShouldBeTrue)
			})
		})

		Convey("->GetObject()", func() {

			Convey("should handle a specific object request", func() {
				json, resp, err := Fetch(baseURL, "test", "1")

				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(json.HasErrors(), ShouldBeFalse)
				So(json.HasData(), ShouldBeTrue)
				So(json.First().ID, ShouldEqual, "1")
			})
		})
	})
}
