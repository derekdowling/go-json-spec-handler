package jsh

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {

	Convey("Parse Tests", t, func() {

		Convey("->validateHeaders()", func() {
			req, reqErr := http.NewRequest("GET", "", nil)
			So(reqErr, ShouldBeNil)
			req.Header.Set("Content-Type", "jpeg")

			err := validateHeaders(req.Header)
			So(err, ShouldNotBeNil)
			So(err.Objects[0].Status, ShouldEqual, http.StatusNotAcceptable)
		})

		Convey("->dataIsArray()", func() {

			Convey("should detect an array successfully", func() {
				data := `[1,2,3,4]`
				isArray, err := dataIsArray([]byte(data))

				So(err, ShouldBeNil)
				So(isArray, ShouldBeTrue)
			})

			Convey("should reject a non-array successfully", func() {
				data := `{"foo":"bar"}`
				isArray, err := dataIsArray([]byte(data))

				So(err, ShouldBeNil)
				So(isArray, ShouldBeFalse)
			})
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
				So(err.Objects[0].Status, ShouldEqual, 422)
				So(err.Objects[0].Source.Pointer, ShouldEqual, "/data/attributes/type")
			})

			Convey("should accept empty ID only for POST", func() {
				objectJSON := `{"data": {"id": "", "type":"test", "attributes": {"ID":"123"}}}`
				req, reqErr := testRequest([]byte(objectJSON))
				So(reqErr, ShouldBeNil)

				Convey("POST test", func() {
					req.Method = "POST"
					_, err := ParseObject(req)
					So(err, ShouldBeNil)
				})

				Convey("PATCH test", func() {
					req.Method = "PATCH"
					_, err := ParseObject(req)
					So(err, ShouldNotBeNil)
				})
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
				So(err.Objects[0].Status, ShouldEqual, 422)
				So(err.Objects[0].Source.Pointer, ShouldEqual, "/data/attributes/id")
			})
		})
	})
}
