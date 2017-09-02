// Copyright (c) Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Client is the REST client type
type Client struct {
	config           *Config
	functionsStackID string
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

// Service refers to rancher's Service
type Service struct {
	ID            string        `json:"id"`
	StackID       string        `json:"stackId"`
	StartOnCreate bool          `json:"startOnCreate"`
	Name          string        `json:"name"`
	Scale         uint64        `json:"scale"`
	LaunchConfig  *LaunchConfig `json:"launchConfig"`
}

func NewClientForConfig(config *Config) (*Client, error) {
	url := fmt.Sprintf("%s/stacks?name=%s", config.CattleURL, config.FunctionsStackName)

	c := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

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

	var objmap map[string]*json.RawMessage
	objmapErr := json.Unmarshal(body, &objmap)
	if objmapErr != nil {
		log.Fatal(objmapErr)
	}

	stacks := make([]Stack, 0)
	jsonErr := json.Unmarshal(*objmap["data"], &stacks)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if len(stacks) == 0 {
		log.Fatal(fmt.Errorf("No stack named %s found", config.FunctionsStackName))
	}
	stack := stacks[0]

	client := Client{
		config:           config,
		functionsStackID: stack.ID,
	}
	return &client, nil
}

func (c *Client) ListServices() ([]Service, error) {
	url := fmt.Sprintf("%s/stacks/%s/services", c.config.CattleURL, c.functionsStackID)

	return c.listServicesInternal(url)
}

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

func (c *Client) listServicesInternal(url string) ([]Service, error) {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(c.config.CattleAccessKey, c.config.CattleSecretKey)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var objmap map[string]*json.RawMessage
	objmapErr := json.Unmarshal(body, &objmap)
	if objmapErr != nil {
		log.Fatal(objmapErr)
	}

	services := make([]Service, 0)
	jsonErr := json.Unmarshal(*objmap["data"], &services)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return services, nil

}

func (c *Client) CreateService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services", c.config.CattleURL)

	spec.StackID = c.functionsStackID
	jsonValue, jsonErr := json.Marshal(spec)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.CattleAccessKey, c.config.CattleSecretKey)

	client := &http.Client{}
	resp, postErr := client.Do(req)
	if postErr != nil {
		log.Fatal(postErr)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, readErr := ioutil.ReadAll(resp.Body)
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

func (c *Client) DeleteService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", c.config.CattleURL, spec.ID)

	req, reqErr := http.NewRequest(http.MethodDelete, url, nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.CattleAccessKey, c.config.CattleSecretKey)

	client := &http.Client{}
	resp, delErr := client.Do(req)
	if delErr != nil {
		log.Fatal(delErr)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, readErr := ioutil.ReadAll(resp.Body)
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

func (c *Client) UpgradeService(spec *Service) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", c.config.CattleURL, spec.ID)

	jsonValue, jsonErr := json.Marshal(spec)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.CattleAccessKey, c.config.CattleSecretKey)

	client := &http.Client{}
	resp, putErr := client.Do(req)
	if putErr != nil {
		log.Fatal(putErr)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, readErr := ioutil.ReadAll(resp.Body)
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
