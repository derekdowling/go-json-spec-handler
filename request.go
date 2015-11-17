package jsh

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// ContentType is the data encoding of choice for HTTP Request and Response Headers
	ContentType = "application/vnd.api+json"
)

// ParseObject returns a JSON object for a given io.ReadCloser containing
// a raw JSON payload. Here's an example of how to use it as part of your full flow.
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		obj, error := jsh.ParseObject(r)
//		if error != nil {
//			// log your error
//			jsh.Send(w, r, error)
//			return
//		}
//
//		yourType := &YourType
//
//		err := object.Unmarshal("yourtype", &YourType)
//		if err != nil {
//			jsh.Send(w, r, err)
//			return
//		}
//
//		yourType.ID = obj.ID
//		// do business logic
//
//		response, err := jsh.NewObject(yourType.ID, "yourtype", &yourType)
//		if err != nil {
//			// log error
//			jsh.Send(w, r, err)
//			return
//		}
//
//		jsh.Send(w, r, response)
//	}
func ParseObject(r *http.Request) (*Object, SendableError) {

	byteData, loadErr := loadJSON(r)
	if loadErr != nil {
		return nil, loadErr
	}

	data := struct {
		Object Object `json:"data"`
	}{}

	err := json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		))
	}

	object := &data.Object
	return object, validateInput(object)
}

// ParseList returns a JSON List for a given io.ReadCloser containing
// a raw JSON payload
func ParseList(r *http.Request) ([]*Object, SendableError) {

	byteData, loadErr := loadJSON(r)
	if loadErr != nil {
		return nil, loadErr
	}

	data := struct {
		List []*Object `json:"data"`
	}{List: []*Object{}}

	err := json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		))
	}

	for _, object := range data.List {
		err := validateInput(object)
		if err != nil {
			return nil, err
		}
	}

	return data.List, nil
}

func loadJSON(r *http.Request) ([]byte, SendableError) {
	defer closeReader(r.Body)

	validationErr := validateRequest(r)
	if validationErr != nil {
		return nil, validationErr
	}

	byteData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Error attempting to read request body: %s", err))
	}

	return byteData, nil
}

func closeReader(reader io.ReadCloser) {
	err := reader.Close()
	if err != nil {
		log.Println("Unabled to close request Body")
	}
}

func validateRequest(r *http.Request) SendableError {

	reqContentType := r.Header.Get("Content-Type")
	if reqContentType != ContentType {
		return SpecificationError(fmt.Sprintf(
			"Expected Content-Type header to be %s, got: %s",
			ContentType,
			reqContentType,
		))
	}

	return nil
}

// NewObjectRequest allows you to create a formatted request that can be used with an
// http.Client or for testing
func NewObjectRequest(method string, baseURL *url.URL, object *Object) (*http.Request, error) {

	switch method {
	case "PATCH":
	case "DELETE":
		baseURL.Path = strings.Join([]string{object.Type, object.ID}, "/")
		break
	case "POST":
		break
	default:
		return nil, SpecificationError(fmt.Sprintf(
			"Cannot use HTTP method ''%s' for a JSON Request", method,
		))
	}

	request, err := http.NewRequest(method, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating new HTTP request: %s", err.Error())
	}

	// use Prepare to generate a payload
	payload, err := object.Prepare(request)
	if err != nil {
		return nil, fmt.Errorf("Error preparing object: %s", err.Error())
	}

	content, jsonErr := json.MarshalIndent(payload, "", "  ")
	if jsonErr != nil {
		return nil, fmt.Errorf("Unable to prepare JSON content: %s", jsonErr)
	}

	request.Header.Add("Content-Type", ContentType)
	request.Header.Set("Content-Length", strconv.Itoa(len(content)))
	request.Body = createIOCloser(content)

	return request, nil
}
