package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPatch(t *testing.T) {

	Convey("Patch Tests", t, func() {
		resource := jshapi.NewMockResource("test", 1, nil)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Patch()", func() {
			object, err := jsh.NewObject("2", "test", nil)
			So(err, ShouldBeNil)

			json, resp, patchErr := Patch(baseURL, object)

			So(patchErr, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(json.HasErrors(), ShouldBeFalse)
			So(json.HasData(), ShouldBeTrue)
		})
	})
}
