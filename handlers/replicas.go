// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/kenfdev/faas-rancher/rancher"
	"github.com/kenfdev/faas-rancher/types"
)

// MakeReplicaUpdater updates desired count of replicas
func MakeReplicaUpdater(client rancher.BridgeClient) VarsHandler {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {

		log.Println("Update replicas")

		functionName := vars["name"]

		req := types.ScaleServiceRequest{}
		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		service, findErr := client.FindServiceByName(functionName)
		if findErr != nil {
			w.WriteHeader(500)
			w.Write([]byte("Unable to lookup function deployment " + functionName))
			log.Println(findErr)
			return
		}

		updates := make(map[string]string)
		updates["scale"] = strconv.FormatInt(req.Replicas, 10)
		_, upgradeErr := client.UpdateService(service, updates)
		if upgradeErr != nil {
			w.WriteHeader(500)
			w.Write([]byte("Unable to update function deployment " + functionName))
			log.Println(upgradeErr)
			return
		}

	}
}

// MakeReplicaReader reads the amount of replicas for a deployment
func MakeReplicaReader(client rancher.BridgeClient) VarsHandler {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {

		log.Println("Read replicas")

		functionName := vars["name"]

		functions, err := getServiceList(client)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		var found *requests.Function
		for _, function := range functions {
			if function.Name == functionName {
				found = &function
				break
			}
		}

		if found == nil {
			w.WriteHeader(404)
			return
		}

		functionBytes, _ := json.Marshal(found)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
