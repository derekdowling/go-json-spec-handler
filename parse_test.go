package jsh

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {

	Convey("Request Tests", t, func() {

		Convey("->ParseObject()", func() {

			Convey("should parse a valid object", func() {

				objectJSON := `{"data": {"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}}`

				closer := createIOCloser([]byte(objectJSON))

				object, err := ParseObject(closer)
				So(err, ShouldBeNil)
				So(object, ShouldNotBeEmpty)
				So(object.Type, ShouldEqual, "user")
				So(object.ID, ShouldEqual, "sweetID123")
				So(object.Attributes, ShouldResemble, json.RawMessage(`{"ID":"123"}`))
			})

			Convey("should reject an object with missing attributes", func() {
				objectJSON := `{"data": {"id": "sweetID123", "attributes": {"ID":"123"}}}`
				closer := createIOCloser([]byte(objectJSON))

				_, err := ParseObject(closer)
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

				closer := createIOCloser([]byte(listJSON))

				list, err := ParseList(closer)
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

				closer := createIOCloser([]byte(listJSON))

				_, err := ParseList(closer)
				So(err, ShouldNotBeNil)

				vErr, ok := err.(*Error)
				So(ok, ShouldBeTrue)
				So(vErr.Status, ShouldEqual, 422)
				So(vErr.Source.Pointer, ShouldEqual, "data/attributes/id")
			})
		})
	})
}
