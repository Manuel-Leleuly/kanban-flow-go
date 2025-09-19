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

func TestCreateUserSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	reqBody := models.UserCreateRequest{
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
		Password:  "newUser@123",
	}

	createUserJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/users", strings.NewReader(string(createUserJson)), "")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.UserResponse
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, reqBody.Email, responseBody.Email)
	assert.Equal(t, reqBody.FirstName, responseBody.FirstName)
	assert.Equal(t, reqBody.LastName, responseBody.LastName)
}

func TestCreateUserFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	// failed because of validation
	reqBody := models.UserCreateRequest{
		FirstName: "First!",
		LastName:  "Last!",
		Email:     "notanemail",
		Password:  "passwordwithincorrectformat",
	}

	createUserJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/users", strings.NewReader(string(createUserJson)), "")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ValidationErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// validation is ordered by the key alphabetically
	// last validation ends with a dot
	assert.Len(t, responseBody.Message, 4)
	assert.Equal(t, "email: must be in email format", responseBody.Message[0])
	assert.Equal(t, "first_name: must only contain alphabet", responseBody.Message[1])
	assert.Equal(t, "last_name: must only contain alphabet", responseBody.Message[2])
	assert.Equal(t, "password: must have at least 1 uppercase letter.", responseBody.Message[3])

	// failed because the email is already used
	reqBody = models.UserCreateRequest{
		FirstName: "Test",
		LastName:  "Test",
		Email:     testhelper.TEST_USER.Email,
		Password:  "CorrectlyF0rmattedP@ssword",
	}

	createUserJson, err = json.Marshal(reqBody)
	assert.Nil(t, err)

	request = testhelper.GetHTTPRequest(http.MethodPost, "/iam/v1/users", strings.NewReader(string(createUserJson)), "")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	var emailResonseBody models.ErrorMessage
	err = json.Unmarshal(body, &emailResonseBody)
	assert.Nil(t, err)

	assert.Equal(t, "email is already used", emailResonseBody.Message)
}
