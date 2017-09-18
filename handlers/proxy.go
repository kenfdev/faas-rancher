// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"io/ioutil"
)

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy(httpDoer HttpDoer, stackName string) VarsHandler {

	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		defer r.Body.Close()

		if r.Method != "POST" {
			// Requests other than POST are not suppored yet
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		service := vars["name"]

		stamp := strconv.FormatInt(time.Now().Unix(), 10)

		defer func(when time.Time) {
			seconds := time.Since(when).Seconds()
			log.Printf("[%s] took %f seconds\n", stamp, seconds)
		}(time.Now())

		watchdogPort := 8080

		requestBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		url := fmt.Sprintf("http://%s.%s:%d/", service, stackName, watchdogPort)

		request, _ := http.NewRequest("POST", url, bytes.NewReader(requestBody))

		copyHeaders(&request.Header, &r.Header)

		defer request.Body.Close()

		response, err := httpDoer.Do(request)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			buf := bytes.NewBufferString("Can't reach service: " + service)
			w.Write(buf.Bytes())
			return
		}

		clientHeader := w.Header()
		copyHeaders(&clientHeader, &response.Header)

		responseBody, _ := ioutil.ReadAll(response.Body)

		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)

	}
}

func copyHeaders(destination *http.Header, source *http.Header) {
	for k, vv := range *source {
		vvClone := make([]string, len(vv))
		copy(vvClone, vv)
		(*destination)[k] = vvClone
	}
}
