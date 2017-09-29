package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kenfdev/faas-rancher/mocks"
	"github.com/kenfdev/faas/gateway/requests"
	"github.com/rancher/go-rancher/v3"
	"github.com/stretchr/testify/assert"
)

func Test_MakeFunctionReader_Get_Service_List_Error(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeFunctionReader(mockClient)

	req, reqErr := http.NewRequest("GET", "/system/functions", nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	mockClient.On("ListServices").Return(nil, fmt.Errorf("Error"))

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusInternalServerError)
	mockClient.AssertExpectations(t)
}

func Test_MakeFunctionReader_Get_Service_List_No_Active_Services(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeFunctionReader(mockClient)

	req, reqErr := http.NewRequest("GET", "/system/functions", nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	nonActiveService := client.Service{
		State: "activating",
	}

	services := []client.Service{
		nonActiveService,
	}
	mockClient.On("ListServices").Return(services, nil)

	// Act
	handler(rr, req, nil)

	// Assert
	responseBody, _ := ioutil.ReadAll(rr.Body)
	responseServices := make([]client.Service, 0)
	json.Unmarshal(responseBody, &responseServices)

	assert.Equal(rr.Code, http.StatusOK)
	assert.Equal(0, len(responseServices))
	mockClient.AssertExpectations(t)
}

func Test_MakeFunctionReader_Get_Service_List_Has_Active_Services(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeFunctionReader(mockClient)

	req, reqErr := http.NewRequest("GET", "/system/functions", nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	nonActiveService := client.Service{
		State: "activating",
	}

	activeService := client.Service{
		State: "active",
		Name:  "SomeFunction",
		Scale: 1,
		LaunchConfig: &client.LaunchConfig{
			ImageUuid: "some/docker/image",
			Labels: map[string]string{
				"faas_function": "some_function",
			},
		},
	}

	replicas := uint64(activeService.Scale)
	expectedFunction := requests.Function{
		Name:            activeService.Name,
		Replicas:        replicas,
		Image:           activeService.LaunchConfig.ImageUuid,
		InvocationCount: 0,
	}

	services := []client.Service{
		nonActiveService,
		activeService,
	}
	mockClient.On("ListServices").Return(services, nil)

	// Act
	handler(rr, req, nil)

	// Assert
	responseBody, _ := ioutil.ReadAll(rr.Body)
	functions := make([]requests.Function, 0)
	json.Unmarshal(responseBody, &functions)

	assert.Equal(rr.Code, http.StatusOK)
	assert.Equal(1, len(functions))

	assert.Equal(expectedFunction, functions[0])
	mockClient.AssertExpectations(t)
}

func Test_MakeFunctionReader_Get_Service_List_Has_Active_Services_But_Not_Labeled(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeFunctionReader(mockClient)

	req, reqErr := http.NewRequest("GET", "/system/functions", nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	nonActiveService := client.Service{
		State: "activating",
	}

	activeButNotLabeledService := client.Service{
		State: "active",
		Name:  "SomeFunction",
		Scale: 1,
		LaunchConfig: &client.LaunchConfig{
			ImageUuid: "some/docker/image",
			// no label to indicate faas function
			Labels: map[string]string{},
		},
	}

	services := []client.Service{
		nonActiveService,
		activeButNotLabeledService,
	}
	mockClient.On("ListServices").Return(services, nil)

	// Act
	handler(rr, req, nil)

	// Assert
	responseBody, _ := ioutil.ReadAll(rr.Body)
	functions := make([]requests.Function, 0)
	json.Unmarshal(responseBody, &functions)

	assert.Equal(rr.Code, http.StatusOK)
	assert.Equal(0, len(functions))

	mockClient.AssertExpectations(t)
}
