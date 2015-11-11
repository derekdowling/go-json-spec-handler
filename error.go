package japi

// Error represents a JSON Spec Error
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}
