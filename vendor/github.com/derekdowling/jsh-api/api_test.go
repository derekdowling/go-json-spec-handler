package jshapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {

	Convey("API Tests", t, func() {

		api := New("foo", nil)

		Convey("->AddResource()", func() {
			resource := NewMockResource("", "test", 1, nil)
			api.Add(resource)

			So(resource.prefix, ShouldEqual, "/foo")
			So(api.Resources["test"], ShouldEqual, resource)
		})
	})
}
