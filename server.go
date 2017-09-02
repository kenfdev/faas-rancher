// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	httpClient := http.Client{
		Timeout: time.Second * TimeoutSeconds,
	}
	rancherClient, err := rancher.NewClientForConfig(config, &httpClient)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()

	r.HandleFunc("/system/functions", handlers.MakeFunctionReader(rancherClient).ServeHTTP).Methods("GET")
	r.HandleFunc("/system/functions", handlers.MakeDeployHandler(rancherClient).ServeHTTP).Methods("POST")
	r.HandleFunc("/system/functions", handlers.MakeDeleteHandler(rancherClient).ServeHTTP).Methods("DELETE")

	r.HandleFunc("/system/function/{name:[-a-zA-Z_0-9]+}", handlers.MakeReplicaReader(rancherClient).ServeHTTP).Methods("GET")
	r.HandleFunc("/system/scale-function/{name:[-a-zA-Z_0-9]+}", handlers.MakeReplicaUpdater(rancherClient).ServeHTTP).Methods("POST")

	functionProxy := handlers.MakeProxy(config.FunctionsStackName).ServeHTTP
	r.HandleFunc("/function/{name:[-a-zA-Z_0-9]+}", functionProxy)
	r.HandleFunc("/function/{name:[-a-zA-Z_0-9]+}/", functionProxy)

	readTimeout := 8 * time.Second
	writeTimeout := 8 * time.Second
	tcpPort := 8080

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", tcpPort),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1MB - can be overridden by setting Server.MaxHeaderBytes.
		Handler:        r,
	}

	log.Fatal(s.ListenAndServe())
}
