package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
