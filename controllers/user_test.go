package controllers_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"temperature/test"
)

func TestUser_Login(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)

	requestBody := strings.NewReader(`{
		"mobile": "***",
		"code": "1234"
		
	}`)
	request, err := http.NewRequest("POST", "/temperature/v1/login", requestBody)
	assert.Nil(err)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(http.StatusOK, response.Code, response)

	body, err := ioutil.ReadAll(response.Body)
	assert.Nil(err)
	fmt.Println(string(body))

}
