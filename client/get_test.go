package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	api := testAPI()
	server := httptest.NewServer(api)
	defer server.Close()

	baseURL := server.URL

	Convey("Get Tests", t, func() {

		Convey("->Get()", func() {
			resp, err := Get(baseURL + "/tests/1")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("->GetList()", func() {

			Convey("should handle an object listing request", func() {
				list, _, err := GetList(baseURL, "test")
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
			})
		})

		Convey("->GetObject()", func() {

			Convey("should handle a specific object request", func() {
				obj, _, err := GetObject(baseURL, "test", "1")
				So(err, ShouldBeNil)
				So(obj.ID, ShouldEqual, "1")
			})
		})
	})
}
