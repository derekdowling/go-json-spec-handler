package jsc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/derekdowling/go-json-spec-handler"
)

// Response is a wrapper around an http.Response that allows us to perform
// intelligent actions on them
type Response struct {
	*http.Response
}

// GetObject validates the HTTP response and parses out the JSON object from the
// body if possible
func (r *Response) GetObject() (*jsh.Object, *jsh.Error) {
	obj, objErr := buildParser(r).GetObject()
	if objErr != nil {
		return nil, objErr.(*jsh.Error)
	}

	return obj, nil
}

// GetList validates the HTTP response and parses out the JSON list from the
// body if possible
func (r *Response) GetList() (jsh.List, *jsh.Error) {
	list, listErr := buildParser(r).GetList()
	if listErr != nil {
		return nil, listErr.(*jsh.Error)
	}

	return list, nil
}

// BodyStr is a convenience function that parses the body of the response into a
// string BUT DOESN'T close the ReadCloser
func (r *Response) BodyStr() (string, *jsh.Error) {

	byteData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", jsh.ISE(fmt.Sprintf("Error attempting to read request body: %s", err.Error()))
	}

	return string(byteData), nil
}

func buildParser(response *Response) *jsh.Parser {
	return &jsh.Parser{
		Method:  "",
		Headers: response.Header,
		Payload: response.Body,
	}
}

// setPath builds a JSON url.Path for a given resource type. Typically this just
// envolves concatting a pluralized resource name
func setPath(url *url.URL, resource string) {

	if url.Path != "" && !strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + "/"
	}

	url.Path = fmt.Sprintf("%s%ss", url.Path, resource)
}

// setIDPath creates a JSON url.Path for a specific resource type including an
// ID specifier.
func setIDPath(url *url.URL, resource string, id string) {
	setPath(url, resource)

	if id != "" {
		url.Path = strings.Join([]string{url.Path, id}, "/")
	}
}

// objectToPayload first prepares/validates the object to ensure it is JSON
// spec compatible, and then marshals it to JSON
func objectToPayload(request *http.Request, object *jsh.Object) ([]byte, *jsh.Error) {

	payload, err := object.Prepare(request, false)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error preparing object: %s", err.Error()))
	}

	jsonContent, jsonErr := json.MarshalIndent(payload, "", "  ")
	if jsonErr != nil {
		return nil, jsh.ISE(fmt.Sprintf("Unable to prepare JSON content: %s", jsonErr))
	}

	return jsonContent, nil
}

// sendPayloadRequest is required for sending JSON payload related requests
// because by default the http package does not set Content-Length headers
func sendObjectRequest(request *http.Request, object *jsh.Object) (*Response, *jsh.Error) {

	payload, err := objectToPayload(request, object)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error converting object to JSON: %s", err.Error()))
	}

	// prepare payload and corresponding headers
	request.Body = jsh.CreateReadCloser(payload)
	request.Header.Add("Content-Type", jsh.ContentType)
	request.Header.Set("Content-Length", strconv.Itoa(len(payload)))

	client := &http.Client{}
	response, clientErr := client.Do(request)
	if clientErr != nil {
		return nil, jsh.ISE(fmt.Sprintf(
			"Error sending %s request: %s", request.Method, clientErr.Error(),
		))
	}

	return &Response{response}, nil
}
