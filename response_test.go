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

					request.Method = "GET"

					err := Send(writer, request, object)
					So(err, ShouldBeNil)
					So(writer.Code, ShouldEqual, http.StatusOK)

					contentLength, convErr := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
					So(convErr, ShouldBeNil)
					So(contentLength, ShouldBeGreaterThan, 0)
					So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)
				})
			})
		})

		Convey("->Ok()", func() {
			ok := Ok()
			err := Send(writer, request, ok)
			So(err, ShouldBeNil)
			So(writer.Code, ShouldEqual, http.StatusOK)
		})
	})
}
