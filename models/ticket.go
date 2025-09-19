package models

import (
	"time"

	"github.com/Manuel-Leleuly/kanban-flow-go/helpers"
	"gorm.io/gorm"
)

type Ticket struct {
	ID          string         `gorm:"column:id;primary_key;not null;<-create" json:"id"`
	Title       string         `gorm:"column:title;not null;" json:"title"`
	Description string         `gorm:"column:description;" json:"decription"`
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

type TicketUpdateRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Assignees   StringArray `json:"assignees"`
	Status      string      `json:"status"`
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
