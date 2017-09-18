// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alexellis/faas-provider"
	bootTypes "github.com/alexellis/faas-provider/types"
	"github.com/kenfdev/faas-rancher/handlers"
	"github.com/kenfdev/faas-rancher/rancher"
)

const (
	// TimeoutSeconds seconds untile timeout for http client
	TimeoutSeconds = 2
)

func main() {
	functionStackName := os.Getenv("FUNCTION_STACK_NAME")
	cattleURL := os.Getenv("CATTLE_URL")
	cattleAccessKey := os.Getenv("CATTLE_ACCESS_KEY")
	cattleSecretKey := os.Getenv("CATTLE_SECRET_KEY")

	// creates the rancher client config
	config, err := rancher.NewClientConfig(
		functionStackName,
		cattleURL,
		cattleAccessKey,
		cattleSecretKey)
	if err != nil {
		panic(err.Error())
	}

	// create the rancher REST client
	rancherClient, err := rancher.NewClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Created Rancher Client")

	proxyClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
	}
	bootstrapHandlers := bootTypes.FaaSHandlers{
		FunctionProxy:  handlers.MakeProxy(&proxyClient, config.FunctionsStackName).ServeHTTP,
		DeleteHandler:  handlers.MakeDeleteHandler(rancherClient).ServeHTTP,
		DeployHandler:  handlers.MakeDeployHandler(rancherClient).ServeHTTP,
		FunctionReader: handlers.MakeFunctionReader(rancherClient).ServeHTTP,
		ReplicaReader:  handlers.MakeReplicaReader(rancherClient).ServeHTTP,
		ReplicaUpdater: handlers.MakeReplicaUpdater(rancherClient).ServeHTTP,
	}
	var port int
	port = 8080
	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:  time.Second * 8,
		WriteTimeout: time.Second * 8,
		TCPPort:      &port,
	}

	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)

}
