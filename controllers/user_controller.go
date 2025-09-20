package controllers

import (
	"net/http"
	"strings"

	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser 		godoc
//
//	@Summary		Create user
//	@Description	Create a user
//	@Tags			User
//	@Router			/iam/v1/users [post]
//	@Accept			json
//	@Produce		json
//	@Param			requestBody	body		models.UserCreateRequest{}	true	"Request Body"
//	@Success		201			{object}	models.UserResponse{}
//	@Failure		400			{object}	models.ErrorMessage{}
//	@Failure		500			{object}	models.ErrorMessage{}
func CreateUser(d *models.DBInstance, c *gin.Context) {
	var reqBody models.UserCreateRequest
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

	var user models.User
	d.DB.Where("email = ?", reqBody.Email).First(&user)
	if user.ID != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "email is already used",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to create user",
		})
		return
	}

	newUser := models.User{
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Email:     reqBody.Email,
		Password:  string(hash),
	}
	if d.DB.Create(&newUser).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, newUser.ToUserResponse())
}
