package models

import (
	"time"

	"github.com/Manuel-Leleuly/kanban-flow-go/helpers"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type Ticket struct {
	ID          string         `gorm:"column:id;primary_key;not null;<-create" json:"id"`
	Title       string         `gorm:"column:title;not null;" json:"title"`
	Description string         `gorm:"column:description;" json:"description"`
	Assignees   StringArray    `gorm:"column:assignees;type:jsonb" json:"assignees"`
	Status      string         `gorm:"column:status;not null;" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime;not null;<-create" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// belongs to
	UserID string `json:"user_id"`
	User   User   `json:"user"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}

func (t *Ticket) BeforeCreate(db *gorm.DB) error {
	if t.ID == "" {
		t.ID = helpers.GenerateUUIDWithoutHyphen()
	}
	return nil
}

func (t *Ticket) ToTicketResponse() TicketResponse {
	return TicketResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Assignees:   t.Assignees,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// request body
type TicketCreateRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Assignees   StringArray `json:"assignees"`
	Status      string      `json:"status"`
}

func (tcr TicketCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&tcr,
		/*
			Title validations:
			- is required
			- min length 8
			- max length 50
		*/
		validation.Field(
			&tcr.Title,
			validation.Required.Error("is required"),
			validation.Length(8, 50).Error("must have length between 8 and 50"),
		),

		/*
			Description validations:
			- min length 1
			- max length 200
		*/
		validation.Field(
			&tcr.Description,
			validation.Length(1, 200).Error("must have length between 1 and 200"),
		),

		/*
			Assignees validations:
			- only allows frontend, backend, and ui
			- must not contain duplicates
		*/
		validation.Field(
			&tcr.Assignees,
			validation.Each(
				validation.In("frontend", "backend", "design").Error("only allows \"frontend\", \"backend\", or \"design\""),
			),
			tcr.Assignees.ValidateUniqueItems(),
		),

		/*
			Status validations:
			- only allows todo, doing, done
		*/
		validation.Field(
			&tcr.Status,
			validation.In("todo", "doing", "done").Error("only allows \"todo\", \"doing\", or \"done\""),
		),
	)
}

type TicketUpdateRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Assignees   StringArray `json:"assignees"`
	Status      string      `json:"status"`
}

func (tur TicketUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&tur,
		/*
			Title validations:
			- is required
			- min length 8
			- max length 50
		*/
		validation.Field(
			&tur.Title,
			validation.Required.Error("is required"),
			validation.Length(8, 50).Error("must have length between 8 and 50"),
		),

		/*
			Description validations:
			- min length 1
			- max length 200
		*/
		validation.Field(
			&tur.Description,
			validation.Length(1, 200).Error("must have length between 1 and 200"),
		),

		/*
			Assignees validations:
			- only allows frontend, backend, and design
			- must not contain duplicates
		*/
		validation.Field(
			&tur.Assignees,
			validation.Each(
				validation.In("frontend", "backend", "design").Error("only allows \"frontend\", \"backend\", or \"design\""),
			),
			tur.Assignees.ValidateUniqueItems(),
		),

		/*
			Status validations:
			- only allows todo, doing, done
		*/
		validation.Field(
			&tur.Status,
			validation.In("todo", "doing", "done").Error("only allows \"todo\", \"doing\", or \"done\""),
		),
	)
}

// response
type TicketResponse struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Assignees   StringArray `json:"assignees"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type TicketDeleteResponse struct {
	Message string `json:"message"`
}
