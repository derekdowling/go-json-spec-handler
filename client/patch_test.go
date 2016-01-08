package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPatch(t *testing.T) {

	Convey("Patch Tests", t, func() {

		api := testAPI()
		server := httptest.NewServer(api)
		defer server.Close()

		baseURL := server.URL

		Convey("->Patch()", func() {
			object, err := jsh.NewObject("2", "tests", nil)
			So(err, ShouldBeNil)

			json, resp, patchErr := Patch(baseURL, object)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(patchErr, ShouldBeNil)
			So(json.HasErrors(), ShouldBeFalse)
			So(json.HasData(), ShouldBeTrue)
		})
	})
}
