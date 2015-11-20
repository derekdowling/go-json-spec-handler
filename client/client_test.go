package jsc

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	. "github.com/smartystreets/goconvey/convey"
)

const testURL = "https://httpbin.org"

func TestClient(t *testing.T) {

	Convey("Client Tests", t, func() {

		Convey("->NewRequest()", func() {

			Convey("should create a valid HTTP request", func() {
				url := &url.URL{Host: "test123"}
				obj := &jsh.Object{ID: "test123", Type: "obj"}
				req, err := NewRequest("POST", url.String(), obj)

				So(err, ShouldBeNil)
				So(req.request.Method, ShouldEqual, "POST")
				So(req.request.URL, ShouldResemble, url)
			})

			Convey("should error for invalid HTTP methods", func() {
				obj := &jsh.Object{}
				_, err := NewRequest("PUT", "", obj)
				So(err, ShouldNotBeNil)

				singleErr, ok := err.(*jsh.Error)
				So(ok, ShouldBeTrue)
				So(singleErr.Status, ShouldEqual, http.StatusNotAcceptable)
			})
		})

		Convey("->Send()", func() {
			obj := &jsh.Object{ID: "test123", Type: "obj"}
			req, err := NewRequest("POST", testURL, obj)
			So(err, ShouldBeNil)

			resp, err := req.Send()
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 405)
		})
	})
}
