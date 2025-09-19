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
		Password:  "newuser123",
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
