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

// LaunchConfig refers to the rancher service's launch config
type LaunchConfig struct {
	Environment   map[string]string `json:"environment"`
	Labels        map[string]string `json:"labels"`
	RestartPolicy map[string]string `json:"restartPolicy"`
	ImageUUID     string            `json:"imageUuid"`
}

// ServiceResponse is the response structure for service requests
type ServiceResponse struct {
	Data []Service `json:"data"`
}

// Service refers to rancher's Service
type Service struct {
	ID            string        `json:"id"`
	StackID       string        `json:"stackId"`
	StartOnCreate bool          `json:"startOnCreate"`
	Name          string        `json:"name"`
	Scale         uint64        `json:"scale"`
	LaunchConfig  *LaunchConfig `json:"launchConfig"`
}

// HTTPClient is a Http Client Wrapper
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
