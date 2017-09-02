// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/kenfdev/faas-rancher/rancher"
)

func getServiceList(client *rancher.Client) ([]requests.Function, error) {
	functions := []requests.Function{}

	services, err := client.ListServices()

	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if _, ok := service.LaunchConfig.Labels["faas_function"]; ok {
			function := requests.Function{
				Name:            service.Name,
				Replicas:        service.Scale,
				Image:           service.LaunchConfig.ImageUUID,
				InvocationCount: 0,
			}
			functions = append(functions, function)

		}
	}

	return functions, nil
}

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(client *rancher.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		functions, err := getServiceList(client)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
