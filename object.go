package japi

// Object represents the default JSON spec for objects
type Object struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
}
