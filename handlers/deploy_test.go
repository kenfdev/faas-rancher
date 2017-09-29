package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/kenfdev/faas-rancher/mocks"
	"github.com/kenfdev/faas/gateway/requests"
	"github.com/rancher/go-rancher/v3"
	"github.com/stretchr/testify/assert"
)

func Test_MakeDeployHandler_Create_Service_Success(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeployHandler(mockClient)

	request := requests.CreateFunctionRequest{
		Service: "some-service",
		EnvVars: map[string]string{
			"SOME_ENV": "SOME_VALUE",
		},
		EnvProcess: "path/to/process",
	}
	b, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	mockClient.On("CreateService",
		mock.MatchedBy(func(s *client.Service) bool {
			return s.Name == request.Service &&
				s.Scale == 1 &&
				s.LaunchConfig.Environment["SOME_ENV"] == request.EnvVars["SOME_ENV"] &&
				s.LaunchConfig.Environment["fprocess"] == request.EnvProcess &&
				s.LaunchConfig.Labels["faas_function"] == request.Service
		}),
	).Return(nil, nil)
	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusAccepted)
}

func Test_MakeDeployHandler_Bad_Json_Request(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeployHandler(mockClient)

	badJSON := []byte(`{name: what?}`)
	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(badJSON))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func Test_MakeDeployHandler_Invalid_Service_Name(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeployHandler(mockClient)

	invalidRequest := requests.CreateFunctionRequest{
		Service: "invalid_servicename", // no valid DNS name
	}
	b, err := json.Marshal(invalidRequest)
	if err != nil {
		log.Fatal(err)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func Test_MakeDeployHandler_Create_Service_Error(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeployHandler(mockClient)

	request := requests.CreateFunctionRequest{
		Service: "some-service",
	}
	b, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	mockClient.On("CreateService",
		mock.MatchedBy(func(s *client.Service) bool { return s.Name == request.Service }),
	).Return(nil, fmt.Errorf("Error"))
	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusInternalServerError)
}
