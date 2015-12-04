package jsc

import (
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/jsh-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	Convey("Get Tests", t, func() {

		resource := jshapi.NewMockResource("", "test", 1, nil)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Get()", func() {

			Convey("should handle an object listing request", func() {
				resp, err := Get(baseURL, "test", "")
				So(err, ShouldBeNil)

				list, err := resp.GetList()
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
			})

			Convey("should handle a specific object request", func() {
				resp, err := Get(baseURL, "test", "1")
				So(err, ShouldBeNil)

				obj, err := resp.GetObject()
				So(err, ShouldBeNil)
				So(obj.ID, ShouldEqual, "1")
			})
		})
	})
}
