package models

import (
	"regexp"
	"time"

	"github.com/Manuel-Leleuly/kanban-flow-go/helpers"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primary_key;column:id;not null;<-create" json:"id"`
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

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = helpers.GenerateUUIDWithoutHyphen()
	}
	return nil
}

func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
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

func (ucr UserCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&ucr,
		/*
			FirstName validations:
			- required
			- min length 2
			- max length 50
			- only contains alphabet
		*/
		validation.Field(
			&ucr.FirstName,
			validation.Required.Error("is required"),
			validation.Length(2, 50).Error("must have length between 2 and 50"),
			is.Alpha.Error("must only contain alphabet"),
		),

		/*
			LastName validations:
			- min length 2
			- max length 50
			- only contains alphabet
		*/
		validation.Field(
			&ucr.LastName,
			validation.Length(2, 50).Error("must have length between 2 and 50"),
			is.Alpha.Error("must only contain alphabet"),
		),

		/*
			Email validations:
			- required
			- must be in email format
		*/
		validation.Field(
			&ucr.Email,
			validation.Required.Error("is required"),
			is.Email.Error("must be in email format"),
		),

		/*
			Password validations:
			- required
			- min length 8
			- max length 50
			- must contain at least 1 uppercase letter
			- must contain at least 1 lowercase letter
			- must contain at least 1 number
			- must contain at least 1 special character
		*/
		validation.Field(&ucr.Password,
			validation.Required.Error("is required"),
			validation.Length(8, 50).Error("must have length between 8 and 50"),
			validation.Match(regexp.MustCompile("[A-Z]+")).Error("must have at least 1 uppercase letter"),
			validation.Match(regexp.MustCompile("[a-z]+")).Error("must have at least 1 lowercase letter"),
			validation.Match(regexp.MustCompile("[0-9]+")).Error("must have at least 1 number"),
			validation.Match(regexp.MustCompile("\\W+")).Error("must have at least 1 non-alphanumeric character"),
		),
	)
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (udr UserUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&udr,
		/*
			FirstName validations:
			- required
			- min length 2
			- max length 50
			- only contains alphabet
		*/
		validation.Field(
			&udr.FirstName,
			validation.Required,
			validation.Length(2, 50),
			is.Alpha,
		),

		/*
			LastName validations:
			- min length 2
			- max length 50
			- only contains alphabet
		*/
		validation.Field(
			&udr.LastName,
			validation.Length(2, 50),
			is.Alpha,
		),
	)
}

// response
type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
