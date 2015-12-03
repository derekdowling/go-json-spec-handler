package jsh

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testURL = "https://httpbin.org"

func TestClientRequest(t *testing.T) {

	Convey("Client Request Tests", t, func() {

		Convey("->setPath()", func() {
			url := &url.URL{Host: "test"}

			Convey("should format properly", func() {
				setPath(url, "test", "1")
				So(url.String(), ShouldEqual, "//test/tests/1")
			})

			Convey("should work with an empty ID", func() {
				setPath(url, "test", "")
				So(url.String(), ShouldEqual, "//test/tests")
			})

			Convey("should respect an existing path", func() {
				url.Path = "admin"
				setPath(url, "test", "")
				So(url.String(), ShouldEqual, "//test/admin/tests")
			})
		})

		Convey("->NewRequest()", func() {

			Convey("should create a valid HTTP request", func() {
				url := &url.URL{Host: "test123"}
				obj := &Object{ID: "123", Type: "obj"}
				req, err := NewRequest("POST", url.String(), obj)

				So(err, ShouldBeNil)
				So(req.Method, ShouldEqual, "POST")
				So(req.URL.String(), ShouldResemble, "//test123/objs/123")
			})

			Convey("should error for invalid HTTP methods", func() {

				Convey("PUT", func() {
					obj := &Object{}
					_, err := NewRequest("PUT", "", obj)
					So(err, ShouldNotBeNil)

					singleErr, ok := err.(*Error)
					So(ok, ShouldBeTrue)
					So(singleErr.Status, ShouldEqual, http.StatusNotAcceptable)
				})

				Convey("GET", func() {
					_, err := NewRequest("GET", "", &Object{})
					So(err, ShouldNotBeNil)

					singleErr, ok := err.(*Error)
					So(ok, ShouldBeTrue)
					So(singleErr.Status, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("->NewGetRequest()", func() {

			Convey("should properly format a base resource Get", func() {
				request, err := NewGetRequest("http://test123", "object", "")
				So(err, ShouldBeNil)

				So(request.URL.Host, ShouldEqual, "test123")
				So(request.URL.Path, ShouldEqual, "/objects")
			})

			Convey("should properly format a specific resource Get", func() {
				request, err := NewGetRequest("http://test123", "object", "345")
				So(err, ShouldBeNil)

				So(request.URL.Host, ShouldEqual, "test123")
				So(request.URL.Path, ShouldEqual, "/objects/345")
			})
		})

		Convey("->Send()", func() {
			obj := &Object{ID: "test123", Type: "obj"}
			req, err := NewRequest("POST", testURL, obj)
			log.Printf("req.URL.String() = %+v\n", req.URL.String())
			So(err, ShouldBeNil)

			resp, err := req.Send()
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 404)
		})
	})
}

func TestClientResponse(t *testing.T) {

	Convey("Client Response Tests", t, func() {

		Convey("->GetObject()", func() {

			obj, objErr := NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)
			r, err := mockClientObjectResponse(obj)
			So(err, ShouldBeNil)
			response := &ClientResponse{r}

			Convey("should parse successfully", func() {
				respObj, err := response.GetObject()
				So(err, ShouldBeNil)
				So(respObj, ShouldNotBeNil)
			})

		})

		Convey("->GetList()", func() {

			obj, objErr := NewObject("123", "test", map[string]string{"test": "test"})
			So(objErr, ShouldBeNil)

			list := &List{}
			list.Add(obj)
			list.Add(obj)

			r, err := mockClientListResponse(list)
			So(err, ShouldBeNil)
			response := &ClientResponse{r}

			Convey("should parse successfully", func() {
				respObj, err := response.GetList()
				So(err, ShouldBeNil)
				So(respObj, ShouldNotBeNil)
			})
		})
	})
}

func mockClientObjectResponse(object *Object) (*http.Response, error) {
	object.ID = "1"

	req, err := NewGetRequest("", object.Type, object.ID)
	if err != nil {
		return nil, err
	}

	resp, err := object.prepare(req.Request, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	SendResponse(recorder, req.Request, resp)
	return recorderToResponse(recorder), nil
}

func mockClientListResponse(list *List) (*http.Response, error) {

	req, err := NewGetRequest("", list.Objects[0].Type, "")
	if err != nil {
		return nil, err
	}

	resp, err := list.prepare(req.Request, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	SendResponse(recorder, req.Request, resp)
	return recorderToResponse(recorder), nil
}

func recorderToResponse(recorder *httptest.ResponseRecorder) *http.Response {
	return &http.Response{
		StatusCode: recorder.Code,
		Body:       CreateReadCloser(recorder.Body.Bytes()),
		Header:     recorder.Header(),
	}
}
