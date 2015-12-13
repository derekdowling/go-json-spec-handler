package jsc

import (
	"net/url"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	. "github.com/smartystreets/goconvey/convey"
)

const testURL = "https://httpbin.org"

func TestClientRequest(t *testing.T) {

	Convey("Client Tests", t, func() {

		Convey("->setPath()", func() {
			url := &url.URL{Host: "test"}

			Convey("should format properly", func() {
				setPath(url, "test")
				So(url.String(), ShouldEqual, "//test/tests")
			})

			Convey("should respect an existing path", func() {
				url.Path = "admin"
				setPath(url, "test")
				So(url.String(), ShouldEqual, "//test/admin/tests")
			})
		})

		Convey("->setIDPath()", func() {
			url := &url.URL{Host: "test"}

			Convey("should format properly an id url", func() {
				setIDPath(url, "test", "1")
				So(url.String(), ShouldEqual, "//test/tests/1")
			})
		})

	})
}

func TestResponseParsing(t *testing.T) {

	Convey("Response Parsing Tests", t, func() {

		Convey("->ParseObject()", func() {

			obj, objErr := jsh.NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)
			response, err := mockObjectResponse(obj)
			So(err, ShouldBeNil)

			Convey("should parse successfully", func() {
				respObj, err := ParseObject(response)
				So(err, ShouldBeNil)
				So(respObj.ID, ShouldEqual, "123")
			})
		})

		Convey("->GetList()", func() {

			obj, objErr := jsh.NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)

			list := jsh.List{obj, obj}

			response, err := mockListResponse(list)
			So(err, ShouldBeNil)

			Convey("should parse successfully", func() {
				list, err := ParseList(response)
				So(err, ShouldBeNil)
				So(list, ShouldNotBeNil)
				So(list[0].ID, ShouldEqual, "123")
			})
		})
	})
}
