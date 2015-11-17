package jsh

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {

	Convey("Request Tests", t, func() {

		Convey("->validateRequest()", func() {
			req, err := http.NewRequest("GET", "", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", "jpeg")

			err = validateRequest(req)
			So(err, ShouldNotBeNil)

			singleErr := err.(*Error)
			So(singleErr.Status, ShouldEqual, http.StatusNotAcceptable)
		})

		Convey("->loadJSON()", func() {
			req, err := http.NewRequest("GET", "", createIOCloser([]byte("1234")))
			So(err, ShouldBeNil)
			req.Header.Set("Content-Type", ContentType)

			bytes, err := loadJSON(req)
			So(err, ShouldBeNil)
			So(string(bytes), ShouldEqual, "1234")
		})

		Convey("->ParseObject()", func() {

			Convey("should parse a valid object", func() {

				objectJSON := `{"data": {"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}}`

				req, reqErr := testRequest([]byte(objectJSON))
				So(reqErr, ShouldBeNil)

				object, err := ParseObject(req)
				So(err, ShouldBeNil)
				So(object, ShouldNotBeEmpty)
				So(object.Type, ShouldEqual, "user")
				So(object.ID, ShouldEqual, "sweetID123")
				So(object.Attributes, ShouldResemble, json.RawMessage(`{"ID":"123"}`))
			})

			Convey("should reject an object with missing attributes", func() {
				objectJSON := `{"data": {"id": "sweetID123", "attributes": {"ID":"123"}}}`

				req, reqErr := testRequest([]byte(objectJSON))
				So(reqErr, ShouldBeNil)

				_, err := ParseObject(req)
				So(err, ShouldNotBeNil)

				vErr, ok := err.(*Error)
				So(ok, ShouldBeTrue)
				So(vErr.Status, ShouldEqual, 422)
				So(vErr.Source.Pointer, ShouldEqual, "data/attributes/type")
			})
		})

		Convey("->ParseList()", func() {

			Convey("should parse a valid list", func() {

				listJSON :=
					`{"data": [
	{"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}},
	{"type": "user", "id": "sweetID456", "attributes": {"ID":"456"}}
]}`
				req, reqErr := testRequest([]byte(listJSON))
				So(reqErr, ShouldBeNil)

				list, err := ParseList(req)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 2)

				object := list[1]
				So(object.Type, ShouldEqual, "user")
				So(object.ID, ShouldEqual, "sweetID456")
				So(object.Attributes, ShouldResemble, json.RawMessage(`{"ID":"456"}`))
			})

			Convey("should error for an invalid list", func() {
				listJSON :=
					`{"data": [
	{"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}},
	{"type": "user", "attributes": {"ID":"456"}}
]}`

				req, reqErr := testRequest([]byte(listJSON))
				So(reqErr, ShouldBeNil)

				_, err := ParseList(req)
				So(err, ShouldNotBeNil)

				vErr, ok := err.(*Error)
				So(ok, ShouldBeTrue)
				So(vErr.Status, ShouldEqual, 422)
				So(vErr.Source.Pointer, ShouldEqual, "data/attributes/id")
			})
		})

		Convey("->NewObjectRequest()", func() {

			Convey("should create a valid HTTP request", func() {
				url := &url.URL{Host: "test123"}
				obj := &Object{ID: "test123", Type: "obj"}
				req, err := NewObjectRequest("POST", url, obj)

				So(err, ShouldBeNil)
				So(req.Method, ShouldEqual, "POST")
				So(req.URL, ShouldResemble, url)
			})

			Convey("should error for invalid HTTP methods", func() {
				url := &url.URL{}
				obj := &Object{}
				_, err := NewObjectRequest("PUT", url, obj)
				So(err, ShouldNotBeNil)

				singleErr, ok := err.(*Error)
				So(ok, ShouldBeTrue)
				So(singleErr.Status, ShouldEqual, http.StatusNotAcceptable)
			})
		})
	})
}
