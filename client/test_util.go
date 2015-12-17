package jsc

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
)

func testAPI() *jshapi.API {
	resource := jshapi.NewMockResource("test", 1, nil)
	api := jshapi.New("", nil)
	api.Add(resource)

	return api
}

func mockObjectResponse(object *jsh.Object) (*http.Response, error) {
	url := &url.URL{Host: "test"}
	setIDPath(url, object.Type, object.ID)

	req, reqErr := http.NewRequest("GET", url.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	resp, err := object.Prepare(req, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	jsh.SendResponse(recorder, req, resp)
	return recorderToResponse(recorder), nil
}

func mockListResponse(list jsh.List) (*http.Response, error) {

	url := &url.URL{Host: "test"}
	setPath(url, list[0].Type)

	req, reqErr := http.NewRequest("GET", url.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	resp, err := list.Prepare(req, false)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	jsh.SendResponse(recorder, req, resp)
	return recorderToResponse(recorder), nil
}

func recorderToResponse(recorder *httptest.ResponseRecorder) *http.Response {
	return &http.Response{
		StatusCode: recorder.Code,
		Body:       jsh.CreateReadCloser(recorder.Body.Bytes()),
		Header:     recorder.Header(),
	}
}
