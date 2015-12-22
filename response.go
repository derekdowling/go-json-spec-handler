package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// JSONAPIVersion is version of JSON API Spec that is currently compatible:
// http://jsonapi.org/format/1.1/
const JSONAPIVersion = "1.1"

// Sendable implements functions that allows different response types
// to produce a sendable JSON Response format
type Sendable interface {
	Validate(r *http.Request, response bool) *Error
}

// Send will return a JSON payload to the requestor. If the payload response validation
// fails, it will send an appropriate error to the requestor and will return the error
func Send(w http.ResponseWriter, r *http.Request, payload Sendable) *Error {

	validationErr := payload.Validate(r, true)
	if validationErr != nil {

		// use the prepared error as the new response, unless something went horribly
		// wrong
		err := validationErr.Validate(r, true)
		if err != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return err
		}

		payload = validationErr
	}

	return SendDocument(w, r, Build(payload))
}

/*
SendDocument handles sending a fully prepared JSON Document. This is useful if you
require custom validation or additional build steps before sending.

SendJSON is designed to always send a response, but will also return the last
error it encountered to help with debugging in the event of an Internal Server
Error.
*/
func SendDocument(w http.ResponseWriter, r *http.Request, document *Document) *Error {

	validationErr := document.Validate(r, true)
	if validationErr != nil {
		prepErr := validationErr.Validate(r, true)

		// If we ever hit this, something seriously wrong has happened
		if prepErr != nil {
			http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
			return prepErr
		}

		// if we didn't error out, make this the new response
		document = Build(validationErr)
	}

	content, jsonErr := json.MarshalIndent(document, "", " ")
	if jsonErr != nil {
		http.Error(w, DefaultErrorTitle, http.StatusInternalServerError)
		return ISE(fmt.Sprintf("Unable to marshal JSON payload: %s", jsonErr.Error()))
	}

	w.Header().Add("Content-Type", ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(document.Status)
	w.Write(content)

	return validationErr
}

// Ok makes it simple to return a 200 OK response via jsh:
//
//	jsh.SendDocument(w, r, jsh.Ok())
func Ok() *Document {
	doc := New()
	doc.Status = http.StatusOK
	doc.empty = true

	return doc
}
