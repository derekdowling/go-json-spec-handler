package jsh

import (
	"encoding/json"
	"log"
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
		req := &http.Request{Method: "GET"}

		Convey("->Validate()", func() {
			err := testList.Validate(req, true)
			So(err, ShouldBeNil)
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

		Convey("->UnmarshalJSON()", func() {

			Convey("should handle a data object", func() {
				jObj := `{"data": {"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}}`

				l := List{}
				err := l.UnmarshalJSON([]byte(jObj))
				log.Printf("l = %+v\n", l)
				So(err, ShouldBeNil)
				So(l, ShouldNotBeEmpty)
			})

			Convey("should handle a data list", func() {
				jList := `{"data": [{"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}]}`

				l := List{}
				err := l.UnmarshalJSON([]byte(jList))
				log.Printf("l = %+v\n", l)
				So(err, ShouldBeNil)
				So(l, ShouldNotBeEmpty)
			})
		})
	})
}
