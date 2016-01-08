package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	api := testAPI()
	server := httptest.NewServer(api)
	defer server.Close()

	baseURL := server.URL

	Convey("Get Tests", t, func() {

		Convey("->Get()", func() {
			json, resp, err := Get(baseURL + "/tests/1")

			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(json.HasErrors(), ShouldBeFalse)
			So(json.HasData(), ShouldBeTrue)
		})

		Convey("->GetList()", func() {

			Convey("should handle an object listing request", func() {
				json, resp, err := List(baseURL, "tests")

				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(json.HasErrors(), ShouldBeFalse)
				So(json.HasData(), ShouldBeTrue)
			})
		})

		Convey("->GetObject()", func() {

			Convey("should handle a specific object request", func() {
				json, resp, err := Fetch(baseURL, "tests", "1")

				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(json.HasErrors(), ShouldBeFalse)
				So(json.HasData(), ShouldBeTrue)
				So(json.First().ID, ShouldEqual, "1")
			})
		})
	})
}
