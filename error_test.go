package jsh

import (
	"log"
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

		testError := &Error{
			Status: http.StatusBadRequest,
			Title:  "Fail",
			Detail: "So badly",
		}

		Convey("->validateError()", func() {

			Convey("should not fail for a valid Error", func() {
				validErr := ISE("Valid Error")
				err := validateError(validErr)
				So(err, ShouldBeNil)
			})

			Convey("422 Status Formatting", func() {

				testError.Status = 422

				Convey("should accept a properly formatted 422 error", func() {
					testError.Source.Pointer = "data/attributes/test"
					err := validateError(testError)
					So(err, ShouldBeNil)
				})

				Convey("should error if Source.Pointer isn't set", func() {
					err := validateError(testError)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("should fail for an out of range HTTP error status", func() {
				testError.Status = http.StatusOK
				err := validateError(testError)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("->Send()", func() {

			testErrors := &ErrorList{Errors: []*Error{&Error{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "Can't Go Here",
			}, testError}}

			Convey("->Send(Error)", func() {
				err := Send(request, writer, testError)
				So(err, ShouldBeNil)
				log.Printf("err = %+v\n", err)
				log.Printf("writer = %+v\n", writer)
				So(writer.Code, ShouldEqual, http.StatusBadRequest)
			})

			Convey("should send a properly formatted JSON Error", func() {
				err := Send(request, writer, testErrors)

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
