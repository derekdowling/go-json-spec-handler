package jsh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestList(t *testing.T) {

	Convey("List Tests", t, func() {

		testObject := &Object{
			ID:         "ID123",
			Type:       "testConversion",
			Attributes: json.RawMessage(`{"foo":"bar"}`),
		}

		testList := List{testObject}
		req := &http.Request{}

		Convey("->Prepare()", func() {
			response, err := testList.Prepare(req, true)
			So(err, ShouldBeNil)
			So(response.HTTPStatus, ShouldEqual, http.StatusOK)
		})

		Convey("->Send(list)", func() {

			Convey("should send a properly formatted List response", func() {

				writer := httptest.NewRecorder()
				err := Send(writer, req, testList)
				So(err, ShouldBeNil)
				So(writer.Code, ShouldEqual, http.StatusOK)

				contentLength, convErr := strconv.Atoi(writer.HeaderMap.Get("Content-Length"))
				So(convErr, ShouldBeNil)
				So(contentLength, ShouldBeGreaterThan, 0)
				So(writer.HeaderMap.Get("Content-Type"), ShouldEqual, ContentType)

				req, reqErr := testRequest(writer.Body.Bytes())
				So(reqErr, ShouldBeNil)

				responseList, parseErr := ParseList(req)
				So(parseErr, ShouldBeNil)
				So(len(responseList), ShouldEqual, 1)
			})
		})
	})
}
