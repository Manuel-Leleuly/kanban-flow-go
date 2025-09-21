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

// TODO: improve this
var accessSecret string = os.Getenv("CLIENT_SECRET") + "_access"
var refreshSecret string = os.Getenv("CLIENT_SECRET") + "_refresh"

func CreateToken(user models.User, isRefresh bool) (string, error) {
	expiredUnix := time.Now().Add(time.Hour).Unix()
	if isRefresh {
		expiredUnix = time.Now().Add(24 * time.Hour).Unix()
	}

	tokenClaims := models.TokenClaims{
		ID:    user.ID,
		Email: user.Email,
		Exp:   expiredUnix,
		Iat:   time.Now().Unix(),
		Jti:   fmt.Sprintf("%d-%d", user.ID, time.Now().UnixNano()),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims.ToJwtClaims())

	secret := accessSecret
	if isRefresh {
		secret = refreshSecret
	}
	return accessToken.SignedString([]byte(secret))
}

func ValidateToken(d *models.DBInstance, tokenString string, isRefresh bool) (*models.User, error) {
	token, err := GetToken(tokenString, isRefresh)
	if err != nil {
		return nil, err
	}

	var user models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenClaims := ConvertJwtClaimsToTokenClaims(claims)

		if float64(time.Now().Unix()) > float64(tokenClaims.Exp) {
			return nil, errors.New("token is expired")
		}

		result := d.DB.Where("id = ? AND email = ?", tokenClaims.ID, tokenClaims.Email).First(&user)
		if result.Error != nil || user.ID == "" {
			return nil, errors.New("unauthorized access")
		}
	} else {
		return nil, errors.New("unauthorized access")
	}

	return &user, nil
}

// helpers
func GetToken(tokenString string, isRefresh bool) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}

		if isRefresh {
			return []byte(refreshSecret), nil
		}
		return []byte(accessSecret), nil
	})
}

func GetTokenStringFromHeader(bearerToken string) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errors.New("invalid bearer token")
	}

	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

	return tokenString, nil
}

func ConvertJwtClaimsToTokenClaims(claims jwt.MapClaims) *models.TokenClaims {
	return &models.TokenClaims{
		ID:    claims["id"].(string),
		Email: claims["email"].(string),
		Exp:   int64(claims["exp"].(float64)),
		Iat:   int64(claims["iat"].(float64)),
		Jti:   claims["jti"].(string),
	}
}
