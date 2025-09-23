package controllers

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/context"
	jwthelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/jwt"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
//
//	@Summary		login
//	@Description	login
//	@Tags			Auth
//	@Router			/iam/v1/login [post]
//	@Accept			json
//	@Produce		json
//	@Param			requestBody	body		models.Login{}	true	"Request Body"
//	@Success		200			{object}	models.Token{}
//	@Failure		400			{object}	models.ErrorMessage{}
//	@Failure		500			{object}	models.ErrorMessage{}
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

	accessTokenString, err := jwthelper.CreateToken(user, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to generate access token",
		})
		return
	}

	refreshTokenString, err := jwthelper.CreateToken(user, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, models.Token{
		Status:       "success",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})
}

// RefreshToken godoc
//
//	@Summary		get new access and refresh token
//	@Description	get new access and refresh token
//	@Security		ApiKeyAuth
//	@Tags			Auth
//	@Router			/iam/v1/token/refresh [post]
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Token{}
//	@Failure		401	{object}	models.ErrorMessage{}
//	@Failure		500	{object}	models.ErrorMessage{}
func RefreshToken(c *gin.Context) {
	user, err := context.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	accessTokenString, err := jwthelper.CreateToken(*user, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to generate access token",
		})
		return
	}

	refreshTokenString, err := jwthelper.CreateToken(*user, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusCreated, models.Token{
		Status:       "success",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})
}

// Logout godoc
//
//	@Summary		logout
//	@Description	logout
//	@Tags			Auth
//	@Router			/iam/v1/logout [post]
//	@Accept			json
//	@Produce		json
//	@Succcess		200 {object} models.LogoutResponse{}
//	@Failure		400	{object}	models.ErrorMessage{}
func Logout(c *gin.Context) {
	err := context.RemoveUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "user already logged out",
		})
		return
	}

	c.JSON(http.StatusOK, models.LogoutResponse{
		Message: "logout success",
	})
}
