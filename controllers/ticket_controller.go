package controllers

import (
	"net/http"
	"strings"

	"github.com/Manuel-Leleuly/kanban-flow-go/context"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

// CreateTicket 	godoc
//
//	@Summary		Create ticket
//	@Description	Create a ticket
//	@Security		ApiKeyAuth
//	@Tags			Ticket
//	@Router			/kanban/v1/tickets [post]
//	@Accept			json
//	@Produce		json
//	@Param			requestBody	body		models.TicketCreateRequest{}	true	"Request Body"
//	@Success		201			{object}	models.TicketResponse{}
//	@Failure		400			{object}	models.ErrorMessage{}
//	@Failure		500			{object}	models.ErrorMessage{}
func CreateTicket(d *models.DBInstance, c *gin.Context) {
	var reqBody models.TicketCreateRequest
	if err := c.Bind(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid request body",
		})
		return
	}

	if err := reqBody.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: strings.Split(err.Error(), "; "),
		})
		return
	}

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	newTicket := models.Ticket{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Assignees:   reqBody.Assignees,
		Status:      reqBody.Status,
		User:        *user,
	}

	if err := d.DB.Create(&newTicket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to create ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, newTicket.ToTicketResponse())

	websocketMessage := models.WSMessage{
		Event:  "created",
		Ticket: newTicket.ToTicketResponse(),
	}
	msg, err := websocketMessage.ToJsonMarshal()
	BroadcastMessage(msg)
}

// GetTicketList 	godoc
//
//	@Summary		Get a list of tickets
//	@Description	Get a list of tickets created by the user stored in the token
//	@Security		ApiKeyAuth
//	@Tags			Ticket
//	@Router			/kanban/v1/tickets [get]
//	@Accept			json
//	@Produce		json
//	@Param			title	query		string	false	"search by ticket title"
//	@Success		200		{object}	[]models.TicketResponse{}
//	@Failure		400		{object}	models.ErrorMessage{}
//	@Failure		401		{object}	models.ErrorMessage{}
//	@Failure		404		{object}	models.ErrorMessage{}
func GetTicketList(d *models.DBInstance, c *gin.Context) {
	/*
		only returns tickets that belong to the user registered
		in the token
	*/

	title := c.Query("title")

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	dbQuery := d.DB.Where("user_id = ?", user.ID)
	if len(title) > 0 {
		dbQuery = dbQuery.Where("LOWER(title) like LOWER(?)", "%"+title+"%")
	}

	var tickets []models.Ticket
	if err := dbQuery.Find(&tickets).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "failed to get all tickets",
		})
		return
	}

	if len(tickets) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "no tickets found",
		})
		return
	}

	var result []models.TicketResponse
	for _, ticket := range tickets {
		result = append(result, ticket.ToTicketResponse())
	}

	c.JSON(http.StatusOK, result)
}

// GetTicketById 	godoc
//
//	@Summary		Get ticket by the ticket ID
//	@Description	Get ticket by the ticket ID
//	@Security		ApiKeyAuth
//	@Tags			Ticket
//	@Router			/kanban/v1/tickets/{ticketId} [get]
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string	true	"Ticket ID"
//	@Success		200			{object}	models.TicketResponse{}
//	@Failure		401			{object}	models.ErrorMessage{}
//	@Failure		404			{object}	models.ErrorMessage{}
func GetTicketById(d *models.DBInstance, c *gin.Context) {
	ticketId := c.Param("ticketId")

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	var ticket models.Ticket
	result := d.DB.Where("user_id = ? AND Tickets.id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "ticket not found",
		})
		return
	}

	c.JSON(http.StatusOK, ticket.ToTicketResponse())
}

// UpdateTicket 	godoc
//
//	@Summary		Update a ticket
//	@Description	Update a ticket
//	@Security		ApiKeyAuth
//	@Tags			Ticket
//	@Router			/kanban/v1/tickets/{ticketId} [put]
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string							true	"Ticket ID"
//	@Param			requestBody	body		models.TicketUpdateRequest{}	true	"Request Body"
//	@Success		200			{object}	models.TicketResponse{}
//	@Failure		401			{object}	models.ErrorMessage{}
//	@Failure		404			{object}	models.ErrorMessage{}
func UpdateTicket(d *models.DBInstance, c *gin.Context) {
	ticketId := c.Param("ticketId")

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	var reqBody models.TicketUpdateRequest
	if err := c.Bind(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid request body",
		})
		return
	}

	if err := reqBody.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: strings.Split(err.Error(), "; "),
		})
		return
	}

	var ticket models.Ticket

	result := d.DB.Where("user_id = ? AND Tickets.id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "ticket not found",
		})
		return
	}

	ticket.Title = reqBody.Title
	ticket.Description = reqBody.Description
	ticket.Assignees = reqBody.Assignees
	ticket.Status = reqBody.Status

	if err := d.DB.Save(&ticket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to update ticket",
		})
		return
	}

	c.JSON(http.StatusOK, ticket.ToTicketResponse())

	websocketMessage := models.WSMessage{
		Event:  "updated",
		Ticket: ticket.ToTicketResponse(),
	}
	msg, err := websocketMessage.ToJsonMarshal()
	BroadcastMessage(msg)
}

// DeleteTicket 	godoc
//
//	@Summary		Delete ticket
//	@Description	Delete a ticket
//	@Security		ApiKeyAuth
//	@Tags			Ticket
//	@Router			/kanban/v1/tickets/{ticketId} [delete]
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string	true	"Ticket ID"
//	@Success		200			{object}	models.TicketDeleteResponse{}
//	@Failure		401			{object}	models.ErrorMessage{}
//	@Failure		404			{object}	models.ErrorMessage{}
func DeleteTicket(d *models.DBInstance, c *gin.Context) {
	ticketId := c.Param("ticketId")

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	var ticket models.Ticket

	result := d.DB.Where("user_id = ? AND Tickets.id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "ticket not found",
		})
		return
	}

	// soft delete ticket
	if err := d.DB.Where("user_id = ? AND Tickets.id = ?", user.ID, ticketId).Delete(&ticket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to delete ticket",
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketDeleteResponse{
		Message: "success",
	})

	websocketMessage := models.WSMessage{
		Event:  "deleted",
		Ticket: models.TicketResponse{},
	}
	msg, err := websocketMessage.ToJsonMarshal()
	BroadcastMessage(msg)
}
