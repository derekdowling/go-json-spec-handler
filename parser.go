package jsh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
ParseObject validates the HTTP request and returns a JSON object for a given
io.ReadCloser containing a raw JSON payload. Here's an example of how to use it
as part of your full flow.

	func Handler(w http.ResponseWriter, r *http.Request) {
		obj, error := jsh.ParseObject(r)
		if error != nil {
			// log your error
			err := jsh.Send(w, r, error)
			return
		}

		yourType := &YourType{}

		err := object.Unmarshal("yourtype", &yourType)
		if err != nil {
			err := jsh.Send(w, r, err)
			return
		}

		yourType.ID = obj.ID
		// do business logic

		err := object.Marshal(yourType)
		if err != nil {
			// log error
			err := jsh.Send(w, r, err)
			return
		}

		err := jsh.Send(w, r, object)
	}
*/
func ParseObject(r *http.Request) (*Object, *Error) {
	document, err := ParseDoc(r, ObjectMode)
	if err != nil {
		return nil, err
	}

	if !document.HasData() {
		return nil, nil
	}

	object := document.First()
	if r.Method != "POST" && object.ID == "" {
		return nil, InputError("Missing mandatory object attribute", "id")
	}

	return object, nil
}

/*
ParseList validates the HTTP request and returns a resulting list of objects
parsed from the request Body. Use just like ParseObject.
*/
func ParseList(r *http.Request) (List, *Error) {
	document, err := ParseDoc(r, ListMode)
	if err != nil {
		return nil, err
	}

	return document.Data, nil
}

// MaxContentLength is 10MB
// https://github.com/golang/go/blob/abb3c0618b658a41bf91a087f1737412e93ff6d9/src/pkg/net/http/request.go#L617
const MaxContentLength int64 = 10 << 20

/*
ParseDoc parses and returns a top level jsh.Document. In most cases, using
"ParseList" or "ParseObject" is preferable.
*/
func ParseDoc(r *http.Request, mode DocumentMode) (*Document, *Error) {
	return NewParser(r).Document(r.Body, mode)
}

// Parser is an abstraction layer that helps to support parsing JSON payload from
// many types of sources, and allows other libraries to leverage this if desired.
type Parser struct {
	Method  string
	Headers http.Header
}

// NewParser creates a parser from an http.Request
func NewParser(request *http.Request) *Parser {
	return &Parser{
		Method:  request.Method,
		Headers: request.Header,
	}
}

/*
Document returns a single JSON data object from the parser. In the process it will
also validate any data objects against the JSON API.
*/
func (p *Parser) Document(payload io.ReadCloser, mode DocumentMode) (*Document, *Error) {
	defer closeReader(payload)

	err := validateHeaders(p.Headers)
	if err != nil {
		return nil, err
	}

	document := &Document{
		Data: List{},
		Mode: mode,
	}

	decodeErr := json.NewDecoder(io.LimitReader(payload, MaxContentLength)).Decode(document)
	if decodeErr != nil {
		return nil, ISE(fmt.Sprintf("Error parsing JSON Document: %s", decodeErr.Error()))
	}

	// If the document has data, validate against specification
	if document.HasData() {
		for _, object := range document.Data {

			// TODO: currently this doesn't really do any user input
			// validation since it is validating against the jsh
			// "Object" type. Figure out how to options pass the
			// corressponding user object struct in to enable this
			// without making the API super clumsy.
			inputErr := validateInput(object)
			if inputErr != nil {
				return nil, inputErr[0]
			}

			// if we have a list, then all resource objects should have IDs, will
			// cross the bridge of bulk creation if and when there is a use case
			if len(document.Data) > 1 && object.ID == "" {
				return nil, InputError("Object without ID present in list", "id")
			}
		}
	}

	return document, nil
}

/*
closeReader is a deferal helper function for closing a reader and logging any errors that might occur after the fact.
*/
func closeReader(reader io.ReadCloser) {
	err := reader.Close()
	if err != nil {
		log.Println("Unable to close request Body")
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
