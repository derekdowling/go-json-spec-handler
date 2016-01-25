package jsc

import (
	"encoding/json"
	"fmt"
	"io"
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
setPath builds a JSON url.Path for a given resource type.
*/
func setPath(url *url.URL, resource string) {

	// ensure that path is "/" terminated before concatting resource
	if url.Path != "" && !strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + "/"
	}

	// don't pluralize resource automagically, JSON API spec is agnostic
	url.Path = fmt.Sprintf("%s%s", url.Path, resource)
}

/*
setIDPath creates a JSON url.Path for a specific resource type including an
ID specifier.
*/
func setIDPath(url *url.URL, resource string, id string) {
	setPath(url, resource)

	// concat "/:id" if not empty
	if id != "" {
		url.Path = strings.Join([]string{url.Path, id}, "/")
	}
}

/*
prepareBody first prepares/validates the object to ensure it is JSON
spec compatible, and then marshals it to JSON, sets the request body and
corresponding attributes.
*/
func prepareBody(request *http.Request, object *jsh.Object) error {

	err := object.Validate(request, false)
	if err != nil {
		return fmt.Errorf("Error preparing object: %s", err.Error())
	}

	doc := jsh.Build(object)

	jsonContent, jsonErr := json.MarshalIndent(doc, "", " ")
	if jsonErr != nil {
		return fmt.Errorf("Unable to prepare JSON content: %s", jsonErr.Error())
	}

	request.Body = jsh.CreateReadCloser(jsonContent)
	request.ContentLength = int64(len(jsonContent))

	return nil
}

/*
Do sends a the specified request to a JSON API compatible endpoint and
returns the resulting JSON Document if possible along with the response,
and any errors that were encountered while sending, or parsing the
JSON Document.

Useful in conjunction with any of the method Request builders or
for times when you want to send a request to a custom endpoint, but would still
like a JSONAPI response.
*/
func Do(request *http.Request) (*jsh.Document, *http.Response, error) {

	client := &http.Client{}
	response, clientErr := client.Do(request)

	if clientErr != nil {
		return nil, nil, fmt.Errorf(
			"Error sending %s request: %s", request.Method, clientErr.Error(),
		)
	}

	doc, parseErr := ParseResponse(response)
	if parseErr != nil {
		return nil, response, fmt.Errorf("Error parsing response: %s", parseErr.Error())
	}

	return doc, response, parseErr
}

/*
ParseResponse handles parsing an HTTP response into a JSON Document if
possible.
*/
func ParseResponse(response *http.Response) (*jsh.Document, error) {

	skipCodes := []int{
		http.StatusNoContent,
		http.StatusNotFound,
	}

	for _, code := range skipCodes {
		if code == response.StatusCode {
			return nil, nil
		}
	}

	document, err := Document(response)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// NewRequest builds a basic request object with the necessary configurations to
// achieve JSON API compatibility
func NewRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", jsh.ContentType)
	request.Header.Set("Content-Length", strconv.Itoa(int(request.ContentLength)))

	return request, err
}
