package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDelete(t *testing.T) {

	Convey("DELETE Tests", t, func() {

		api := testAPI()
		server := httptest.NewServer(api)
		defer server.Close()

		baseURL := server.URL

		Convey("->Delete()", func() {
			resp, err := Delete(baseURL, "tests", "1")

			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
		})
	})
}
