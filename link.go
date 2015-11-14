package jsh

// Links is a top-level document field
type Links struct {
	Self    *Link `json:"self,omitempty"`
	Related *Link `json:"related,omitempty"`
}

// Link is a JSON format type
type Link struct {
	HREF string                 `json:"href,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}
