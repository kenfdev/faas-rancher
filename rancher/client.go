// Copyright (c) Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package rancher

import (
	"fmt"

	"github.com/rancher/go-rancher/v2"
)

// BridgeClient is the interface for Rancher API
type BridgeClient interface {
	ListServices() ([]client.Service, error)
	FindServiceByName(name string) (*client.Service, error)
	CreateService(spec *client.Service) (*client.Service, error)
	DeleteService(spec *client.Service) error
	UpdateService(spec *client.Service, updates map[string]string) (*client.Service, error)
}

// Client is the REST client type
type Client struct {
	rancherClient    *client.RancherClient
	config           *Config
	functionsStackID string
}

// NewClientForConfig creates a new rancher REST client
func NewClientForConfig(config *Config) (BridgeClient, error) {
	c, newErr := client.NewRancherClient(&client.ClientOpts{
		Url:       config.CattleURL,
		AccessKey: config.CattleAccessKey,
		SecretKey: config.CattleSecretKey,
	})

	if newErr != nil {
		return nil, newErr
	}

	coll, listErr := c.Stack.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"name": config.FunctionsStackName,
		},
	})

	if listErr != nil {
		return nil, listErr
	}

	var stack *client.Stack
	if len(coll.Data) == 0 {
		fmt.Println("stack named " + config.FunctionsStackName + " not found. creating...")
		// create stack if not present
		reqStack := &client.Stack{
			Name: config.FunctionsStackName,
		}
		newStack, err := c.Stack.Create(reqStack)
		if err != nil {
			return nil, err
		}
		fmt.Println("stack creation complete")
		stack = newStack
	} else {
		stack = &coll.Data[0]
	}

	client := Client{
		rancherClient:    c,
		config:           config,
		functionsStackID: stack.Id,
	}

	return &client, nil

}

// ListServices lists rancher services inside the specified stack (set in config)
func (c *Client) ListServices() ([]client.Service, error) {
	services, err := c.rancherClient.Service.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"stackId": c.functionsStackID,
		},
	})
	if err != nil {
		return nil, err
	}
	return services.Data, nil
}

// FindServiceByName finds a service based on its name
func (c *Client) FindServiceByName(name string) (*client.Service, error) {
	services, err := c.rancherClient.Service.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(services.Data) == 0 {
		return nil, fmt.Errorf("No service named " + name + " found.")
	}
	return &services.Data[0], nil
}

// CreateService creates a service inside rancher
func (c *Client) CreateService(spec *client.Service) (*client.Service, error) {

	spec.StackId = c.functionsStackID
	service, err := c.rancherClient.Service.Create(spec)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// DeleteService deletes the specified service in rancher
func (c *Client) DeleteService(spec *client.Service) error {
	err := c.rancherClient.Service.Delete(spec)
	if err != nil {
		return err
	}

	return nil
}

// UpdateService upgrades the specified service in rancher
func (c *Client) UpdateService(spec *client.Service, updates map[string]string) (*client.Service, error) {
	service, err := c.rancherClient.Service.Update(spec, updates)
	if err != nil {
		return nil, err
	}
	return service, nil
}
