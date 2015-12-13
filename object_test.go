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
			Type:       "testObject",
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

			Convey("non-govalidator structs", func() {

				testConversion := struct {
					ID  string
					Foo string `json:"foo"`
				}{}

				Convey("Should successfully populate a valid struct", func() {
					err := testObject.Unmarshal("testObject", &testConversion)
					So(err, ShouldBeNil)
					So(testConversion.Foo, ShouldEqual, "bar")
				})

				Convey("Should reject a non-matching type", func() {
					err := testObject.Unmarshal("badType", &testConversion)
					So(err, ShouldNotBeNil)
				})

			})

			Convey("govalidator struct unmarshals", func() {

				Convey("should not error if input validates properly", func() {
					testValidation := struct {
						Foo string `json:"foo" valid:"alphanum"`
					}{}

					err := testObject.Unmarshal("testObject", &testValidation)
					So(err, ShouldBeNil)
					So(testValidation.Foo, ShouldEqual, "bar")
				})

				Convey("should return a 422 Error correctly for a validation failure", func() {
					testValidation := struct {
						Foo string `valid:"ipv4,required" json:"foo"`
					}{}

					err := testObject.Unmarshal("testObject", &testValidation)
					So(err, ShouldNotBeNil)
					So(err.Objects[0].Source.Pointer, ShouldEqual, "/data/attributes/foo")
				})

				Convey("should return a 422 Error correctly for multiple validation failures", func() {

					testManyObject := &Object{
						ID:         "ID123",
						Type:       "testObject",
						Attributes: json.RawMessage(`{"foo":"bar", "baz":"4567"}`),
					}

					testManyValidations := struct {
						Foo string `valid:"ipv4,required" json:"foo"`
						Baz string `valid:"alpha,required" json:"baz"`
					}{}

					err := testManyObject.Unmarshal("testObject", &testManyValidations)
					So(err, ShouldNotBeNil)

					So(err.Objects[0].Source.Pointer, ShouldEqual, "/data/attributes/foo")
					So(err.Objects[1].Source.Pointer, ShouldEqual, "/data/attributes/baz")
				})
			})
		})

		Convey("->Marshal()", func() {

			Convey("should properly update attributes", func() {
				attrs := map[string]string{"foo": "baz"}
				err := testObject.Marshal(attrs)
				So(err, ShouldBeNil)

				raw, jsonErr := json.MarshalIndent(attrs, "", " ")
				So(jsonErr, ShouldBeNil)
				So(string(testObject.Attributes), ShouldEqual, string(raw))
			})
		})

		Convey("->Prepare()", func() {

			Convey("should handle a POST response correctly", func() {
				request.Method = "POST"
				resp, err := testObject.Prepare(request, true)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusCreated)
			})

			Convey("should handle a GET response correctly", func() {
				request.Method = "GET"
				resp, err := testObject.Prepare(request, true)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusOK)
			})

			Convey("should handle a PATCH response correctly", func() {
				request.Method = "PATCH"
				resp, err := testObject.Prepare(request, true)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusOK)
			})

			Convey("should return a formatted Error for an unsupported method Type", func() {
				request.Method = "PUT"
				resp, err := testObject.Prepare(request, true)
				So(err, ShouldBeNil)
				So(resp.HTTPStatus, ShouldEqual, http.StatusNotAcceptable)
			})
		})

		Convey("->Send(Object)", func() {
			request.Method = "POST"
			writer := httptest.NewRecorder()
			err := Send(writer, request, testObject)
			So(err, ShouldBeNil)
			So(writer.Code, ShouldEqual, http.StatusCreated)
		})
	})
}
