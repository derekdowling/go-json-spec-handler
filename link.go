package jsh

import "encoding/json"

// Links is a top-level document field
type Links struct {
	Self    *Link `json:"self,omitempty"`
	Related *Link `json:"related,omitempty"`
}

// Link is a resource link that can encode as a string or as an object
// as per the JSON API specification.
type Link struct {
	link
}

type link struct {
	HREF string                 `json:"href,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// NewLink creates a new link encoded as a string.
func NewLink(href string) *Link {
	return &Link{
		link: link{
			HREF: href,
		},
	}
}

// NewMetaLink creates a new link with metadata encoded as an object.
func NewMetaLink(href string, meta map[string]interface{}) *Link {
	return &Link{
		link: link{
			HREF: href,
			Meta: meta,
		},
	}
}

// MarshalJSON implements the Marshaler interface for Link.
func (l *Link) MarshalJSON() ([]byte, error) {
	if l.Meta == nil {
		return json.Marshal(l.HREF)
	}
	return json.Marshal(l.link)
}

// UnmarshalJSON implements the Unmarshaler interface for Link.
func (l *Link) UnmarshalJSON(data []byte) error {
	var href string
	err := json.Unmarshal(data, &href)
	if err == nil {
		l.HREF = href
		return nil
	}
	link := link{}
	err = json.Unmarshal(data, &link)
	if err != nil {
		return err
	}
	l.link = link
	return nil
}
