package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	Id          uuid.UUID      `gorm:"column:id;primary_key;not null;type:uuid;<-create" json:"id"`
	Title       string         `gorm:"column:title;not null;" json:"title"`
	Description string         `gorm:"column:description;" json:"decription"`
	Assignees   StringArray    `gorm:"column:assignees;type:jsonb" json:"assignees"`
	Status      string         `gorm:"column:status;not null;" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime;not null;<-create" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	// belongs to
	UserID uuid.UUID
	User   User
}

func (t *Ticket) TableName() string {
	return "tickets"
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
