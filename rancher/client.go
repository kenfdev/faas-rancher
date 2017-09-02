// Copyright (c) Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// NewClientForConfig creates a new rancher REST client
func NewClientForConfig(config *Config, c HTTPClient) (*Client, error) {
	url := fmt.Sprintf("%s/stacks?name=%s", config.CattleURL, config.FunctionsStackName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.CattleAccessKey, config.CattleSecretKey)

	res, getErr := c.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var stackResp StackResponse
	jsonErr := json.Unmarshal(body, &stackResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if len(stackResp.Data) == 0 {
		log.Fatal(fmt.Errorf("No stack named %s found", config.FunctionsStackName))
	}
	stack := stackResp.Data[0]

	client := Client{
		http:             c,
		config:           config,
		functionsStackID: stack.ID,
	}
	return &client, nil
}

// ListServices lists rancher services inside the specified stack (set in config)
func (c *Client) ListServices() ([]Service, error) {
	url := fmt.Sprintf("%s/stacks/%s/services", c.config.CattleURL, c.functionsStackID)

	return c.listServicesInternal(url)
}

// FindServiceByName finds a service based on its name
func (c *Client) FindServiceByName(name string) (*Service, error) {
	url := fmt.Sprintf("%s/services?name=%s", c.config.CattleURL, name)
	services, err := c.listServicesInternal(url)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no service %s found", name)
	}

	return &services[0], nil
}

func (c *Client) execute(method, urlStr string, body io.Reader) ([]byte, error) {
	req, reqErr := http.NewRequest(method, urlStr, body)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.CattleAccessKey, c.config.CattleSecretKey)

	resp, doErr := c.http.Do(req)
	if doErr != nil {
		return nil, doErr
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return ioutil.ReadAll(resp.Body)

}

func (c *Client) listServicesInternal(url string) ([]Service, error) {

	body, exeErr := c.execute(http.MethodGet, url, nil)
	if exeErr != nil {
		log.Fatal(exeErr)
	}

	var serviceResp ServiceResponse
	jsonErr := json.Unmarshal(body, &serviceResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return serviceResp.Data, nil

}

// CreateService creates a service inside rancher
func (c *Client) CreateService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services", c.config.CattleURL)

	spec.StackID = c.functionsStackID
	jsonValue, jsonErr := json.Marshal(spec)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// TODO what happens if create service that already exists
	body, exeErr := c.execute(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if exeErr != nil {
		log.Fatal(exeErr)
	}
	fmt.Println("response Body:", string(body))

	service := Service{}
	unmarshalErr := json.Unmarshal(body, &service)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}
	return &service, nil
}

// DeleteService deletes the specified service in rancher
func (c *Client) DeleteService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", c.config.CattleURL, spec.ID)

	body, readErr := c.execute(http.MethodDelete, url, nil)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println("response Body:", string(body))

	service := Service{}
	unmarshalErr := json.Unmarshal(body, &service)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}
	return &service, nil
}

// UpgradeService upgrades the specified service in rancher
func (c *Client) UpgradeService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", c.config.CattleURL, spec.ID)

	jsonValue, jsonErr := json.Marshal(spec)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	body, readErr := c.execute(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println("response Body:", string(body))

	service := Service{}
	unmarshalErr := json.Unmarshal(body, &service)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}
	return &service, nil
}
