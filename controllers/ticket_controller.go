package controllers

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/context"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func CreateTicket(d *models.DBInstance, c *gin.Context) {
	var reqBody models.TicketCreateRequest
	if err := c.Bind(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid request body",
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

	c.JSON(http.StatusOK, newTicket.ToTicketResponse())
}

func GetTicketList(d *models.DBInstance, c *gin.Context) {
	/*
		only returns tickets that belong to the user registered
		in the token
	*/

	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	var tickets []models.Ticket
	if err := d.DB.Where("user_id = ?", user.ID).Find(&tickets).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "failed to get all tickets",
		})
		return
	}

	var result []models.TicketResponse
	for _, ticket := range tickets {
		result = append(result, ticket.ToTicketResponse())
	}

	c.JSON(http.StatusOK, result)
}

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
	result := d.DB.Where("user_id = ? AND id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID.String() == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "ticket not found",
		})
		return
	}

	c.JSON(http.StatusOK, ticket.ToTicketResponse())
}

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

	var ticket models.Ticket

	result := d.DB.Where("user_id = ? AND id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID.String() == "" {
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
}

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

	result := d.DB.Where("user_id = ? AND id = ?", user.ID, ticketId).First(&ticket)
	if result.Error != nil || ticket.ID.String() == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorMessage{
			Message: "ticket not found",
		})
		return
	}

	// soft delete ticket
	if err := d.DB.Where("user_id = ? AND id = ?", user.ID, ticketId).Delete(&ticket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to delete ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
