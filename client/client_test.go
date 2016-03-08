package jsc

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
	. "github.com/smartystreets/goconvey/convey"
)

const testURL = "https://httpbin.org"

func TestClientRequest(t *testing.T) {

	Convey("Client Tests", t, func() {

		Convey("->setPath()", func() {
			url := &url.URL{Host: "test"}

			Convey("should format properly", func() {
				setPath(url, "tests")
				So(url.String(), ShouldEqual, "//test/tests")
			})

			Convey("should respect an existing path", func() {
				url.Path = "admin"
				setPath(url, "test")
				So(url.String(), ShouldEqual, "//test/admin/test")
			})
		})

		Convey("->setIDPath()", func() {
			url := &url.URL{Host: "test"}

			Convey("should format properly an id url", func() {
				setIDPath(url, "tests", "1")
				So(url.String(), ShouldEqual, "//test/tests/1")
			})
		})

	})
}

func TestParseResponse(t *testing.T) {

	Convey("ParseResponse", t, func() {

		response := &http.Response{
			StatusCode: http.StatusNotFound,
		}

		Convey("404 response parsing should not return a 406 error", func() {
			doc, err := ParseResponse(response)
			So(doc, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestResponseParsing(t *testing.T) {

	Convey("Response Parsing Tests", t, func() {

		Convey("Parse Object", func() {

			obj, objErr := jsh.NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)

			response, err := mockObjectResponse(obj)
			So(err, ShouldBeNil)

			Convey("should parse successfully", func() {
				doc, err := Document(response)

				So(err, ShouldBeNil)
				So(doc.HasData(), ShouldBeTrue)
				So(doc.First().ID, ShouldEqual, "123")
			})
		})

		Convey("Parse List", func() {

			obj, objErr := jsh.NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)

			list := jsh.List{obj, obj}

			response, err := mockListResponse(list)
			So(err, ShouldBeNil)

			Convey("should parse successfully", func() {
				doc, err := Document(response)

				So(err, ShouldBeNil)
				So(doc.HasData(), ShouldBeTrue)
				So(doc.First().ID, ShouldEqual, "123")
			})
		})
	})
}

// not a great for this, would much rather have it in test_util, but it causes an
// import cycle wit jsh-api
func testAPI() *jshapi.API {

	resource := jshapi.NewMockResource("tests", 1, nil)
	resource.Action("testAction", func(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
		object, err := jsh.NewObject("1", "tests", []string{"testAction"})
		if err != nil {
			log.Fatal(err.Error())
		}

		return object, nil
	})

	api := jshapi.New("")
	api.Add(resource)

	return api
}
