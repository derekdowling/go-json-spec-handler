package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAction(t *testing.T) {

	Convey("Action Tests", t, func() {

		api := testAPI()
		server := httptest.NewServer(api)
		defer server.Close()

		baseURL := server.URL

		Convey("->Action()", func() {
			doc, resp, err := Action(baseURL, "tests", "1", "testAction")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(doc, ShouldNotBeEmpty)
		})
	})
}
