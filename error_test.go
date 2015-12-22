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

		request := &http.Request{}
		writer := httptest.NewRecorder()

		testErrorObject := &Error{
			Status: http.StatusBadRequest,
			Title:  "Fail",
			Detail: "So badly",
		}

		Convey("->Validate()", func() {

			Convey("should not fail for a valid Error", func() {
				err := testErrorObject.Validate(request, true)
				So(err, ShouldBeNil)
			})

			Convey("422 Status Formatting", func() {

				testErrorObject.Status = 422

				Convey("should accept a properly formatted 422 error", func() {
					testErrorObject.Source.Pointer = "/data/attributes/test"
					err := testErrorObject.Validate(request, true)
					So(err, ShouldBeNil)
				})

				Convey("should error if Source.Pointer isn't set", func() {
					err := testErrorObject.Validate(request, true)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("should fail for an out of range HTTP error status", func() {
				testErrorObject.Status = http.StatusOK
				err := testErrorObject.Validate(request, true)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("->Send()", func() {

			testError := &Error{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "Can't Go Here",
			}

			Convey("should send a properly formatted JSON error", func() {
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
}
