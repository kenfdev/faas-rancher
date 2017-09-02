// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/kenfdev/faas-rancher/rancher"
)

// ValidateDeployRequest validates that the service name is valid for Kubernetes
func ValidateDeployRequest(request *requests.CreateFunctionRequest) error {
	var validDNS = regexp.MustCompile(`^[a-zA-Z\-]+$`)
	matched := validDNS.MatchString(request.Service)
	if matched {
		return nil
	}

	return fmt.Errorf("(%s) must be a valid DNS entry for service name", request.Service)
}

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler(client *rancher.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		serviceSpec := makeServiceSpec(request)

		_, err = client.CreateService(serviceSpec)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("Created service - " + request.Service)
		log.Println(string(body))

		w.WriteHeader(http.StatusAccepted)

	}
}

func makeServiceSpec(request requests.CreateFunctionRequest) *rancher.Service {
	envVars := request.EnvVars
	if envVars == nil {
		envVars = make(map[string]string)
	}
	envVars["fprocess"] = request.EnvProcess

	restartPolicy := make(map[string]string)
	restartPolicy["name"] = "always"

	labels := make(map[string]string)
	labels["faas_function"] = request.Service
	labels["io.rancher.container.pull_image"] = "always"

	launchConfig := &rancher.LaunchConfig{
		Environment:   envVars,
		RestartPolicy: restartPolicy,
		ImageUUID:     "docker:" + request.Image, // not sure if it's ok to just prefix with 'docker:'
		Labels:        labels,
	}
	serviceSpec := &rancher.Service{
		Name:          request.Service,
		Scale:         1,
		StartOnCreate: true,
		LaunchConfig:  launchConfig,
	}

	return serviceSpec
}
