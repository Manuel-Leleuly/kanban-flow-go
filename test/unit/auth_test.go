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

func TestRefreshTokenSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	tokenData, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/token/refresh", nil, tokenData.RefreshToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.Token
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "success", responseBody.Status)
}

func TestRefreshTokenFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	refreshToken := "wrong_refresh_token"
	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/token/refresh", nil, refreshToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "unauthorized access", responseBody.Message)
}

func TestGetMeSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	token, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodGet, "/iam/v1/users/me", nil, token.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.UserResponse
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, responseBody.ID, testhelper.TEST_USER.ID)
	assert.Equal(t, responseBody.FirstName, testhelper.TEST_USER.FirstName)
	assert.Equal(t, responseBody.LastName, testhelper.TEST_USER.LastName)
	assert.Equal(t, responseBody.Email, testhelper.TEST_USER.Email)
}
