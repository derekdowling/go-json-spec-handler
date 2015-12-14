package jsh

// JSON represents a top level JSON formatted Document.
// Refer to the JSON API Specification for a full descriptor
// of each attribute: http://jsonapi.org/format/#document-structure
type JSON struct {
	Data     *Data       `json:"data,omitempty"`
	Errors   *Error      `json:"errors,omitempty"`
	Links    *Link       `json:"links,omitempty"`
	Included *List       `json:"included,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
	JSONAPI  struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	// HTTPStatus attribute to make it simpler to define expected responses
	// for a given context
	HTTPStatus int `json:"-"`
	// empty is used to signify that the response shouldn't contain a json payload
	// in the case that we only want to return an HTTP Status Code in order to bypass
	// validation steps
	empty bool
}

// Validate checks JSON Spec for the top level JSON document
func (j *JSON) Validate() *Error {

	if !j.empty {
		if !j.HasErrors() && !j.HasData() {
			return ISE("Both `errors` and `data` cannot be blank for a JSON response")
		}
		if j.HasErrors() && j.HasData() {
			return ISE("Both `errors` and `data` cannot be set for a JSON response")
		}
		if j.HasData() && j.Included != nil {
			return ISE("'included' should only be set for a response if 'data' is as well")
		}
	}
	if j.HTTPStatus < 100 || j.HTTPStatus > 600 {
		return ISE("Response HTTP Status is outside of valid range")
	}

	// probably not the best place for this, but...
	j.JSONAPI.Version = JSONAPIVersion

	return nil
}

// HasData will return true if the JSON document's Data field is set
func (j *JSON) HasData() bool {
	return j.Data != nil
}

// HasErrors will return true if the Errors attribute is not nil.
func (j *JSON) HasErrors() bool {
	return j.Errors != nil
}

// Data is an overarching type that deals with parsing json objects and lists
type Data struct {
	Objects []*Object
}
