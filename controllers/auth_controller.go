package controllers

import (
	"net/http"

	jwthelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/jwt"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(d *models.DBInstance, c *gin.Context) {
	var reqBody models.Login
	if err := c.Bind(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid email and/or password",
		})
		return
	}

	var user models.User
	result := d.DB.Where("email = ?", reqBody.Email).First(&user)
	if result.Error != nil || user.ID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid email and/or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "invalid email and/or password",
		})
		return
	}

	accessTokenString, err := jwthelper.CreateAccessToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to generate access token",
		})
		return
	}

	c.Set("me", user)

	c.JSON(http.StatusOK, models.Token{
		Status:      "success",
		AccessToken: accessTokenString,
	})
}
