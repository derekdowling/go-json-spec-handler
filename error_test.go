package jsh

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {

	Convey("Error Tests", t, func() {

		writer := httptest.NewRecorder()
		request := &http.Request{}

		testErrorObject := &ErrorObject{
			Status: http.StatusBadRequest,
			Title:  "Fail",
			Detail: "So badly",
		}

		Convey("->validateError()", func() {

			Convey("should not fail for a valid Error", func() {
				err := validateError(testErrorObject)
				So(err, ShouldBeNil)
			})

			Convey("422 Status Formatting", func() {

				testErrorObject.Status = 422

				Convey("should accept a properly formatted 422 error", func() {
					testErrorObject.Source.Pointer = "/data/attributes/test"
					err := validateError(testErrorObject)
					So(err, ShouldBeNil)
				})

				Convey("should error if Source.Pointer isn't set", func() {
					err := validateError(testErrorObject)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("should fail for an out of range HTTP error status", func() {
				testErrorObject.Status = http.StatusOK
				err := validateError(testErrorObject)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Error List Tests", func() {

			Convey("->Add()", func() {

				testError := &Error{}

				Convey("should successfully add a valid error", func() {
					err := testError.Add(testErrorObject)
					So(err, ShouldBeNil)
					So(len(testError.Objects), ShouldEqual, 1)
				})

				Convey("should error if validation fails while adding an error", func() {
					badError := &ErrorObject{
						Title:  "Invalid",
						Detail: "So badly",
					}

					err := testError.Add(badError)
					So(err.Objects[0].Status, ShouldEqual, 500)
					So(testError.Objects, ShouldBeEmpty)
				})
			})

			Convey("->Send()", func() {

				testError := NewError(&ErrorObject{
					Status: http.StatusForbidden,
					Title:  "Forbidden",
					Detail: "Can't Go Here",
				})

				Convey("should send a properly formatted JSON error list", func() {
					err := Send(writer, request, testError)
					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusForbidden)

					contentLength, convErr := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
					So(convErr, ShouldBeNil)
					So(contentLength, ShouldBeGreaterThan, 0)
					So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)
				})
			})
		})
	})
}
