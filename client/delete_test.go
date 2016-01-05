package jsc

import (
	"log"
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
			resp, err := Delete(baseURL, "test", "1")

			So(err, ShouldBeNil)
			log.Printf("patchErr = %+v\n", err)
			log.Printf("resp = %+v\n", resp)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})
	})
}
