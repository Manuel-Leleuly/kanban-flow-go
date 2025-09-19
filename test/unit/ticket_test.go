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

func TestCreateTicketSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	token, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	reqBody := models.TicketCreateRequest{
		Title:       "New Test Ticket",
		Description: "New Test Ticket Description",
		Assignees:   []string{"frontend"},
		Status:      "todo",
	}

	ticketJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/kanban/v1/tickets", strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.TicketResponse
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, reqBody.Title, responseBody.Title)
	assert.Equal(t, reqBody.Description, responseBody.Description)
	assert.Equal(t, reqBody.Status, responseBody.Status)

	assert.Len(t, responseBody.Assignees, len(reqBody.Assignees))
	assert.Equal(t, reqBody.Assignees[0], responseBody.Assignees[0])
}

func TestCreateTicketFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	token, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	reqBody := models.TicketCreateRequest{
		Title:       "a",
		Description: "",
		Assignees:   []string{"frontend", "frontend"},
		Status:      "sleeping",
	}

	ticketJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPost, "/kanban/v1/tickets", strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ValidationErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// validation key is ordered alphabetically
	// last validation ends with a dot
	assert.Equal(t, "assignees: contains duplicate value frontend", responseBody.Message[0])
	assert.Equal(t, "status: only allows \"todo\", \"doing\", or \"done\"", responseBody.Message[1])
	assert.Equal(t, "title: must have length between 8 and 50.", responseBody.Message[2])
}

func TestUpdateTicketSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	token, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	// update title
	titleReqBody := models.TicketUpdateRequest{
		Title:       "Updated Title",
		Description: testhelper.TEST_TICKET.Description,
		Assignees:   testhelper.TEST_TICKET.Assignees,
		Status:      testhelper.TEST_TICKET.Status,
	}

	ticketJson, err := json.Marshal(titleReqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPut, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var titleResponseBody models.TicketResponse
	err = json.Unmarshal(body, &titleResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, titleReqBody.Title, titleResponseBody.Title)
	assert.Equal(t, titleReqBody.Description, titleResponseBody.Description)
	assert.Equal(t, titleReqBody.Status, titleResponseBody.Status)

	assert.Len(t, titleResponseBody.Assignees, len(titleReqBody.Assignees))
	assert.Equal(t, titleReqBody.Assignees[0], titleResponseBody.Assignees[0])

	// update description
	descriptionReqBody := models.TicketUpdateRequest{
		Title:       titleReqBody.Title,
		Description: "Updated Description",
		Assignees:   testhelper.TEST_TICKET.Assignees,
		Status:      testhelper.TEST_TICKET.Status,
	}

	ticketJson, err = json.Marshal(descriptionReqBody)
	assert.Nil(t, err)

	request = testhelper.GetHTTPRequest(http.MethodPut, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	var descriptionResponseBody models.TicketResponse
	err = json.Unmarshal(body, &descriptionResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, descriptionReqBody.Title, descriptionResponseBody.Title)
	assert.Equal(t, descriptionReqBody.Description, descriptionResponseBody.Description)
	assert.Equal(t, descriptionReqBody.Status, descriptionResponseBody.Status)

	assert.Len(t, descriptionResponseBody.Assignees, len(descriptionReqBody.Assignees))
	assert.Equal(t, descriptionReqBody.Assignees[0], descriptionResponseBody.Assignees[0])

	// update assignees
	assigneesReqBody := models.TicketUpdateRequest{
		Title:       descriptionReqBody.Title,
		Description: descriptionReqBody.Description,
		Assignees:   []string{"frontend", "backend"},
		Status:      testhelper.TEST_TICKET.Status,
	}

	ticketJson, err = json.Marshal(assigneesReqBody)
	assert.Nil(t, err)

	request = testhelper.GetHTTPRequest(http.MethodPut, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	var assigneesResponseBody models.TicketResponse
	err = json.Unmarshal(body, &assigneesResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, assigneesReqBody.Title, assigneesResponseBody.Title)
	assert.Equal(t, assigneesReqBody.Description, assigneesResponseBody.Description)
	assert.Equal(t, assigneesReqBody.Status, assigneesResponseBody.Status)

	assert.Len(t, assigneesResponseBody.Assignees, len(assigneesReqBody.Assignees))
	assert.Equal(t, assigneesReqBody.Assignees[0], assigneesResponseBody.Assignees[0])
	assert.Equal(t, assigneesReqBody.Assignees[1], assigneesResponseBody.Assignees[1])

	// update status
	statusReqBody := models.TicketUpdateRequest{
		Title:       assigneesReqBody.Title,
		Description: assigneesReqBody.Description,
		Assignees:   assigneesReqBody.Assignees,
		Status:      "doing",
	}

	ticketJson, err = json.Marshal(statusReqBody)
	assert.Nil(t, err)

	request = testhelper.GetHTTPRequest(http.MethodPut, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, strings.NewReader(string(ticketJson)), token.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	var statusResponseBody models.TicketResponse
	err = json.Unmarshal(body, &statusResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, statusReqBody.Title, statusResponseBody.Title)
	assert.Equal(t, statusReqBody.Description, statusResponseBody.Description)
	assert.Equal(t, statusReqBody.Status, statusResponseBody.Status)

	assert.Len(t, statusResponseBody.Assignees, len(statusReqBody.Assignees))
	assert.Equal(t, statusReqBody.Assignees[0], statusResponseBody.Assignees[0])
	assert.Equal(t, statusReqBody.Assignees[1], statusResponseBody.Assignees[1])
}

func TestUpdateTicketFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	tokenData, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	// failed because of validation
	reqBody := models.TicketUpdateRequest{
		Title:       "a",
		Description: "",
		Assignees:   []string{"frontend", "frontend"},
		Status:      "sleeping",
	}

	ticketJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodPut, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, strings.NewReader(string(ticketJson)), tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ValidationErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// validation key is ordered alphabetically
	// last validation ends with a dot
	assert.Equal(t, "assignees: contains duplicate value frontend", responseBody.Message[0])
	assert.Equal(t, "status: only allows \"todo\", \"doing\", or \"done\"", responseBody.Message[1])
	assert.Equal(t, "title: must have length between 8 and 50.", responseBody.Message[2])
}

func TestDeleteTicketSuccess(t *testing.T) {
	router := routes.GetRoutes(D)

	tokenData, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodDelete, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, nil, tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.TicketDeleteResponse
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "success", responseBody.Message)

	// check if the ticket is actually deleted (soft delete)
	request = testhelper.GetHTTPRequest(http.MethodGet, "/kanban/v1/tickets/"+testhelper.TEST_TICKET.ID, nil, tokenData.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	var ticketResponseBody models.ErrorMessage
	err = json.Unmarshal(body, &ticketResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, "ticket not found", ticketResponseBody.Message)
}

func TestDeleteTicketFailed(t *testing.T) {
	router := routes.GetRoutes(D)

	tokenData, err := testhelper.GetTestToken(D)
	assert.Nil(t, err)

	request := testhelper.GetHTTPRequest(http.MethodDelete, "/kanban/v1/tickets/wrongticketid", nil, tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "ticket not found", responseBody.Message)
}
