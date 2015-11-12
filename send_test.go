package jsh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSend(t *testing.T) {

	Convey("Send Tests", t, func() {

		writer := httptest.NewRecorder()
		request := &http.Request{}

		object := &Object{
			ID:         "1234",
			Type:       "user",
			Attributes: json.RawMessage(`{"foo":"bar"}`),
		}

		Convey("Success Handlers", func() {

			Convey("->Send()", func() {

				Convey("should send a proper HTTP JSON response", func() {
					err := Send(writer, http.StatusOK, object)
					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusOK)

					contentLength, err := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
					So(err, ShouldBeNil)
					So(contentLength, ShouldBeGreaterThan, 0)
					So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)
				})
			})

			Convey("->SendObject()", func() {

				Convey("should handle a POST response correctly", func() {
					request.Method = "POST"
					err := SendObject(writer, request, object)

					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusCreated)
				})

				Convey("should handle a GET response correctly", func() {
					request.Method = "GET"
					err := SendObject(writer, request, object)

					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusOK)
				})

				Convey("should handle a PATCH response correctly", func() {
					request.Method = "PATCH"
					err := SendObject(writer, request, object)

					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("->SendList()", func() {

				Convey("should send a properly formatted multi-object response", func() {
					list := []*Object{object, object}
					err := SendList(writer, list)

					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusOK)
					contentLength, err := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
					So(err, ShouldBeNil)
					So(contentLength, ShouldBeGreaterThan, 0)
					So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)

					closer := createIOCloser(writer.Body.Bytes())
					responseList, err := ParseList(closer)
					So(len(responseList), ShouldEqual, 2)
				})
			})
		})

		Convey("Error Handlers", func() {
			testError := &Error{
				Status: http.StatusBadRequest,
				Title:  "Fail",
				Detail: "So badly",
			}

			Convey("->SendErrors()", func() {

				testErrors := []*Error{&Error{
					Status: http.StatusForbidden,
					Title:  "Forbidden",
					Detail: "Can't Go Here",
				}, testError}

				Convey("should send a properly formatted JSON Errors", func() {
					err := SendErrors(writer, testErrors)

					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusForbidden)

					contentLength, err := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
					So(err, ShouldBeNil)
					So(contentLength, ShouldBeGreaterThan, 0)
					So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)
				})

				Convey("422 Status Formatting", func() {

					testError.Status = 422

					Convey("should accept a properly formatted 422 error", func() {
						testError.Source.Pointer = "data/attributes/test"
						err := SendError(writer, testError)
						So(err, ShouldBeNil)
					})

					Convey("should reject if err.Source.Pointer is missing", func() {
						err := SendError(writer, testError)
						So(err, ShouldNotBeNil)
					})
				})

				Convey("should reject bad error Status", func() {
					testError.Status = http.StatusOK
					err := SendErrors(writer, testErrors)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("->SendError()", func() {
				err := SendError(writer, testError)
				So(err, ShouldBeNil)
				So(writer.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}
