package unit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	testhelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/test"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/Manuel-Leleuly/kanban-flow-go/routes"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	reqBody := models.Login{
		Email:    testhelper.TEST_USER.Email,
		Password: testhelper.TEST_USER.Password,
	}

	loginJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/login", strings.NewReader(string(loginJson)), "")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.Token
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "success", responseBody.Status)
}

func TestLoginFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	// wrong email
	reqBody := models.Login{
		Email:    "wrong_credentials@example.com",
		Password: testhelper.TEST_USER.Password,
	}

	loginJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/login", strings.NewReader(string(loginJson)), "")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "invalid email and/or password", responseBody.Message)

	// wrong password
	reqBody = models.Login{
		Email:    testhelper.TEST_USER.Email,
		Password: "wrongpassword12345",
	}

	loginJson, err = json.Marshal(reqBody)
	assert.Nil(t, err)

	request = httptest.NewRequest(http.MethodPost, "/iam/v1/login", strings.NewReader(string(loginJson)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "invalid email and/or password", responseBody.Message)
}
