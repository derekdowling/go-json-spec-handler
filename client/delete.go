package jsc

import (
	"fmt"
	"net/http"
	"net/url"
)

// Delete allows a user to make an outbound DELETE /resources/:id request:
//
//	resp, err := jsh.Delete("http://apiserver", "user", "2")
//
func Delete(urlStr string, resourceType string, id string) (*http.Response, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// ghetto pluralization, fix when it becomes an issue
	setIDPath(u, resourceType, id)

	request, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building DELETE request: %s", err.Error())
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error sending DELETE request: %s", err.Error())
	}

	return response, nil
}
