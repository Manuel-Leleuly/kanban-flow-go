package models

type Token struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}
