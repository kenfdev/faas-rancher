// Copyright (c) 2017 Ken Fukuyama
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package rancher

import "net/http"

// Client is the REST client type
type Client struct {
	http             HTTPClient
	config           *Config
	functionsStackID string
}

// StackResponse is the response structure for stack requests
type StackResponse struct {
	Data []Stack `json:"data"`
}

// Stack refers to rancher's stack
type Stack struct {
	ID string `json:"id"`
}

// HTTPClient is a Http Client Wrapper
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
