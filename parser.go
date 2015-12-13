package jsh

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// ContentType is the data encoding of choice for HTTP Request and Response Headers
	ContentType = "application/vnd.api+json"
)

// ParseObject validates the HTTP request and returns a JSON object for a given
// io.ReadCloser containing a raw JSON payload. Here's an example of how to use it
// as part of your full flow.
//
//	func Handler(w http.ResponseWriter, r *http.Request) {
//		obj, error := jsh.ParseObject(r)
//		if error != nil {
//			// log your error
//			err := jsh.Send(w, r, error)
//			return
//		}
//
//		yourType := &YourType{}
//
//		err := object.Unmarshal("yourtype", &yourType)
//		if err != nil {
//			err := jsh.Send(w, r, err)
//			return
//		}
//
//		yourType.ID = obj.ID
//		// do business logic
//
//		err := object.Marshal(yourType)
//		if err != nil {
//			// log error
//			err := jsh.Send(w, r, err)
//			return
//		}
//
//		err := jsh.Send(w, r, object)
//	}
func ParseObject(r *http.Request) (*Object, *Error) {

	object, err := buildParser(r).GetObject()
	if err != nil {
		return nil, err
	}

	if r.Method != "POST" && object.ID == "" {
		return nil, InputError("id", "Missing mandatory object attribute")
	}

	return object, nil
}

// ParseList validates the HTTP request and returns a resulting list of objects
// parsed from the request Body. Use just like ParseObject.
func ParseList(r *http.Request) (List, *Error) {
	return buildParser(r).GetList()
}

// Parser is an abstraction layer to support parsing JSON payload from many types
// of sources in order to allow other packages to use this parser
type Parser struct {
	Method  string
	Headers http.Header
	Payload io.ReadCloser
}

// BuildParser creates a parser from an http.Request
func buildParser(request *http.Request) *Parser {
	return &Parser{
		Method:  request.Method,
		Headers: request.Header,
		Payload: request.Body,
	}
}

// GetObject returns a single JSON data object from the parser
func (p *Parser) GetObject() (*Object, *Error) {
	byteData, loadErr := prepareJSON(p.Headers, p.Payload)
	if loadErr != nil {
		return nil, loadErr
	}

	data := struct {
		Object *Object `json:"data"`
	}{}

	err := json.Unmarshal(byteData, &data)
	if err != nil {
		return nil, ISE(fmt.Sprintf("Unable to parse json: \n%s\nError:%s",
			string(byteData),
			err.Error(),
		))
	}

	object := data.Object

	inputErr := validateInput(object)
	if inputErr != nil {
		return nil, inputErr
	}

	return object, nil
}

// GetList returns a JSON data list from the parser
func (p *Parser) GetList() (List, *Error) {
	byteData, loadErr := prepareJSON(p.Headers, p.Payload)
	if loadErr != nil {
		return nil, loadErr
	}

	data := struct {
		List List `json:"data"`
	}{List{}}

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

		if object.ID == "" {
			return nil, InputError("id", "Object without ID present in list")
		}
	}

	return data.List, nil
}

// prepareJSON ensures that the provide headers are JSON API compatible and then
// reads and closes the closer
func prepareJSON(headers http.Header, closer io.ReadCloser) ([]byte, *Error) {
	defer closeReader(closer)

	validationErr := validateHeaders(headers)
	if validationErr != nil {
		return nil, validationErr
	}

	byteData, err := ioutil.ReadAll(closer)
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

func validateHeaders(headers http.Header) *Error {

	reqContentType := headers.Get("Content-Type")
	if reqContentType != ContentType {
		return SpecificationError(fmt.Sprintf(
			"Expected Content-Type header to be %s, got: %s",
			ContentType,
			reqContentType,
		))
	}

	return nil
}
