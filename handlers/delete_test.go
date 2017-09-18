package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kenfdev/faas-rancher/mocks"
	"github.com/kenfdev/faas/gateway/requests"
	"github.com/rancher/go-rancher/v2"
	"github.com/stretchr/testify/assert"
)

func Test_MakeDeleteHandler_Service_Delete_Success(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)
	functionName := "some_function"

	body := requests.DeleteFunctionRequest{
		FunctionName: functionName,
	}
	b, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	expectedService := client.Service{
		Name: "some_rancher_service",
	}
	mockClient.On("FindServiceByName", functionName).Return(&expectedService, nil)
	mockClient.On("DeleteService", &expectedService).Return(nil)

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusOK)
	mockClient.AssertExpectations(t)
}

func Test_MakeDeleteHandler_InvalidBody(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)

	b := []byte(`{"name":what?}`)
	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func Test_MakeDeleteHandler_Empty_FunctionName(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)
	emptyFunctionName := ""

	invalidBody := requests.DeleteFunctionRequest{
		FunctionName: emptyFunctionName,
	}
	b, jsonErr := json.Marshal(invalidBody)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func Test_MakeDeleteHandler_Service_Find_Error(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)
	functionName := "some_function"

	body := requests.DeleteFunctionRequest{
		FunctionName: functionName,
	}
	b, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	mockClient.On("FindServiceByName", functionName).Return(nil, fmt.Errorf("Internal Server Error"))

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(http.StatusInternalServerError, rr.Code)
	mockClient.AssertExpectations(t)
}

func Test_MakeDeleteHandler_Service_Nil_Error(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)
	functionName := "some_function"

	body := requests.DeleteFunctionRequest{
		FunctionName: functionName,
	}
	b, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	mockClient.On("FindServiceByName", functionName).Return(nil, nil)

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(http.StatusNotFound, rr.Code)
	mockClient.AssertExpectations(t)
}
func Test_MakeDeleteHandler_Service_Delete_Error(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.BridgeClient)
	handler := MakeDeleteHandler(mockClient)
	functionName := "some_function"

	body := requests.DeleteFunctionRequest{
		FunctionName: functionName,
	}
	b, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "/system/functions", bytes.NewReader(b))
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()

	expectedService := client.Service{
		Name: "some_rancher_service",
	}
	mockClient.On("FindServiceByName", functionName).Return(&expectedService, nil)
	mockClient.On("DeleteService", &expectedService).Return(fmt.Errorf("Service Delete Failed"))

	// Act
	handler(rr, req, nil)

	// Assert
	assert.Equal(rr.Code, http.StatusBadRequest)
	mockClient.AssertExpectations(t)
}
