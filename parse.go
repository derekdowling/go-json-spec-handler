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

// ParseObject returns a JSON object for a given io.ReadCloser containing
// a raw JSON payload
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
//		err := jsh.Send(w, r, response)
//		if err != nil {
//			http.Error(w, err.Status, err.Detail)
//		}
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
