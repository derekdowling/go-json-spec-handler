package jsh

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObject(t *testing.T) {

	Convey("Object Tests", t, func() {

		testObject := &Object{
			ID:         "ID123",
			Type:       "testConversion",
			Attributes: json.RawMessage(`{"foo":"bar"}`),
		}

		Convey("->NewObject()", func() {

			Convey("should create a new object with populated attrs", func() {
				attrs := struct {
					Foo string `json:"foo"`
				}{"bar"}

				newObj, err := NewObject(testObject.ID, testObject.Type, attrs)
				So(err, ShouldBeNil)
				So(newObj.Attributes, ShouldNotBeEmpty)
			})
		})

		Convey("->Unmarshal()", func() {
			testConversion := struct {
				ID  string
				Foo string `json:"foo"`
			}{}

			Convey("Should successfully populate a valid struct", func() {
				err := testObject.Unmarshal("testConversion", &testConversion)
				So(err, ShouldBeNil)
				So(testConversion.Foo, ShouldEqual, "bar")
			})

			Convey("Should reject a non-matching type", func() {
				err := testObject.Unmarshal("badType", &testConversion)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("->ParseList()", func() {
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
	})
}
