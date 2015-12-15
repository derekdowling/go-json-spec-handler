package jsc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPost(t *testing.T) {

	attrs := map[string]string{
		"foo": "bar",
	}

	mock := jshapi.NewMockResource("test", 0, attrs)
	server := httptest.NewServer(mock)
	baseURL := server.URL

	Convey("Post Tests", t, func() {
		testObject, err := jsh.NewObject("", "test", attrs)
		So(err, ShouldBeNil)

		_, resp, postErr := Post(baseURL, testObject)
		So(postErr, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, http.StatusCreated)
	})
}
