package jwthelper

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user models.User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	return accessToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
}

func ValidateAccessToken(d *models.DBInstance, accessToken string) error {
	token, err := GetToken(accessToken)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return errors.New("token is expired")
		}

		var user models.User
		result := d.DB.Where("id = ? AND email = ?", claims["id"], claims["email"]).First(&user)
		if result.Error != nil || user.ID.String() == "" {
			return errors.New("unauthorized access")
		}
	} else {
		return errors.New("unauthorized access")
	}

	return nil
}

// helpers
func GetToken(tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("CLIENT_SECRET")), nil
	})
}

func GetTokenStringFromHeader(bearerToken string) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errors.New("invalid bearer token")
	}

	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

	return tokenString, nil
}
