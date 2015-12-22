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

/*
Document validates the HTTP response and attempts to parse a JSON API compatible
Document from the response body before closing it.
*/
func Document(response *http.Response) (*jsh.Document, *jsh.Error) {
	document, err := buildParser(response).Document(response.Body)
	if err != nil {
		return nil, err
	}

	document.Status = response.StatusCode
	return document, nil
}

/*
DumpBody is a convenience function that parses the body of the response into a
string BUT DOESN'T close the ReadCloser. Useful for debugging.
*/
func DumpBody(response *http.Response) (string, *jsh.Error) {

	byteData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", jsh.ISE(fmt.Sprintf("Error attempting to read request body: %s", err.Error()))
	}

	return string(byteData), nil
}

func buildParser(response *http.Response) *jsh.Parser {
	return &jsh.Parser{
		Method:  "",
		Headers: response.Header,
	}
}

/*
setPath builds a JSON url.Path for a given resource type. Typically this just
envolves concatting a pluralized resource name.
*/
func setPath(url *url.URL, resource string) {

	if url.Path != "" && !strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + "/"
	}

	url.Path = fmt.Sprintf("%s%ss", url.Path, resource)
}

/*
setIDPath creates a JSON url.Path for a specific resource type including an
ID specifier.
*/
func setIDPath(url *url.URL, resource string, id string) {
	setPath(url, resource)

	if id != "" {
		url.Path = strings.Join([]string{url.Path, id}, "/")
	}
}

// objectToPayload first prepares/validates the object to ensure it is JSON
// spec compatible, and then marshals it to JSON
func objectToPayload(request *http.Request, object *jsh.Object) ([]byte, *jsh.Error) {

	err := object.Validate(request, false)
	if err != nil {
		return nil, jsh.ISE(fmt.Sprintf("Error preparing object: %s", err.Error()))
	}

	doc := jsh.Build(object)

	jsonContent, jsonErr := json.MarshalIndent(doc, "", "  ")
	if jsonErr != nil {
		return nil, jsh.ISE(fmt.Sprintf("Unable to prepare JSON content: %s", jsonErr))
	}

	return jsonContent, nil
}

// sendPayloadRequest is required for sending JSON payload related requests
// because by default the http package does not set Content-Length headers
func doObjectRequest(request *http.Request, object *jsh.Object) (*jsh.Document, *http.Response, *jsh.Error) {

	payload, err := objectToPayload(request, object)
	if err != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf("Error converting object to JSON: %s", err.Error()))
	}

	// prepare payload and corresponding headers
	request.Body = jsh.CreateReadCloser(payload)
	request.Header.Add("Content-Type", jsh.ContentType)

	contentLength := strconv.Itoa(len(payload))
	request.ContentLength = int64(len(payload))
	request.Header.Set("Content-Length", contentLength)

	client := &http.Client{}
	httpResponse, clientErr := client.Do(request)

	if clientErr != nil {
		return nil, nil, jsh.ISE(fmt.Sprintf(
			"Error sending %s request: %s", request.Method, clientErr.Error(),
		))
	}

	document, err := Document(httpResponse)
	if err != nil {
		return nil, httpResponse, err
	}

	return document, httpResponse, nil
}
