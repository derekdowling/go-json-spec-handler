package jsh

// Error represents a JSON Specification Error.
// Error.Source.Pointer is used in 422 status responses to indicate validation
// errors on a JSON Object attribute.
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}
