package jsc

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

func mockObjectResponse(object *jsh.Object) (*Response, error) {
	object.ID = "1"

	url := &url.URL{Host: "test"}
	setIDPath(url, object.Type, object.ID)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := object.Prepare(req, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	jsh.SendResponse(recorder, req, resp)
	return recorderToResponse(recorder), nil
}

func mockListResponse(list jsh.List) (*Response, error) {

	url := &url.URL{Host: "test"}
	setPath(url, list[0].Type)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := list.Prepare(req, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	jsh.SendResponse(recorder, req, resp)
	return recorderToResponse(recorder), nil
}

func recorderToResponse(recorder *httptest.ResponseRecorder) *Response {
	return &Response{&http.Response{
		StatusCode: recorder.Code,
		Body:       jsh.CreateReadCloser(recorder.Body.Bytes()),
		Header:     recorder.Header(),
	}}
}
