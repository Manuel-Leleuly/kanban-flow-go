package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"primary_key;column:id;not null;<-create" json:"id"`
	FirstName string         `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string         `gorm:"column:last_name;not null" json:"last_name"`
	Email     string         `gorm:"column:email;not null" json:"email"`
	Password  string         `gorm:"column:password;not null" json:"password"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;not null;<-create" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// request body
type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// response
type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
