package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Jti   string `json:"jti"`
}

func (t *TokenClaims) ToJwtClaims() jwt.Claims {
	return jwt.MapClaims{
		"id":    t.ID,
		"email": t.Email,
		"exp":   t.Exp,
		"iat":   t.Iat,
		"jti":   t.Jti,
	}
}
