package handlers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/kenfdev/faas-rancher/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_MakeProxyHandler_Create_Service_Success(t *testing.T) {
	assert := assert.New(t)
	// Arrange
	mockClient := new(mocks.HttpDoer)
	stackName := "some_stackname"
	serviceName := "some-service"
	vars := map[string]string{
		"name": serviceName,
	}
	handler := MakeProxy(mockClient, stackName)

	reqBody := []byte(`{ "data": "some_data" }`)
	req, err := http.NewRequest("POST", "/system/function/"+serviceName, bytes.NewReader(reqBody))
	// req.Header.Add("Content-Type", "Some-Content-Type")
	if err != nil {
		log.Fatal(err)
	}

	responseBody := []byte(`{ "data": "some-data"}`)
	response := &http.Response{
		Header: make(http.Header, 0),
		Body:   ioutil.NopCloser(bytes.NewReader(responseBody)),
	}
	response.Header.Add("Content-Type", "Some-Content-Type")

	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		pBody, _ := ioutil.ReadAll(r.Body)
		// chek if proxied request is called
		return r.Method == req.Method &&
			bytes.Equal(pBody, reqBody)
	})).Return(response, nil)

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req, vars)

	// Assert
	assert.Equal(rr.Code, http.StatusOK)

	proxiedBody, _ := ioutil.ReadAll(rr.Body)
	assert.True(bytes.Equal(responseBody, proxiedBody))
	assert.Equal("Some-Content-Type", rr.Header().Get("Content-Type"), "Headers weren't copied")
}
