package jsh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

		request := &http.Request{}

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

			Convey("with input validation", func() {

				Convey("should not error if input validates properly", func() {

				})

				Convey("should return a 422 Error correctly for validation failure", func() {

				})
			})
		})

		Convey("->Prepare()", func() {

			Convey("should handle a POST response correctly", func() {
				request.Method = "POST"
				resp, err := testObject.Prepare(request)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusCreated)
			})

			Convey("should handle a GET response correctly", func() {
				request.Method = "GET"
				resp, err := testObject.Prepare(request)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusOK)
			})

			Convey("should handle a PATCH response correctly", func() {
				request.Method = "PATCH"
				resp, err := testObject.Prepare(request)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusOK)
			})

			Convey("should return a formatted Error for an unsupported method Type", func() {
				request.Method = "PUT"
				resp, err := testObject.Prepare(request)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusNotAcceptable)
			})
		})

		Convey("->Send(Object)", func() {
			request.Method = "POST"
			writer := httptest.NewRecorder()
			err := Send(request, writer, testObject)
			So(err, ShouldBeNil)
			So(writer.Code, ShouldEqual, http.StatusCreated)
		})
	})
}
