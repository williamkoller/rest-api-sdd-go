package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type TicketModel struct {
	ID          string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID    string  `gorm:"type:uuid;not null;index"`
	UnitID      *string `gorm:"type:uuid"`
	RequesterID string  `gorm:"type:uuid;not null;index"`
	Category    string  `gorm:"not null;size:30"`
	Status      string  `gorm:"not null;default:'open'"`
	Subject     string  `gorm:"not null;size:500"`
	ResolvedAt  *time.Time
	CreatedAt   time.Time            `gorm:"autoCreateTime"`
	UpdatedAt   time.Time            `gorm:"autoUpdateTime"`
	Messages    []TicketMessageModel `gorm:"foreignKey:TicketID"`
}

func (TicketModel) TableName() string { return "tickets" }

func (m *TicketModel) ToEntity() *entity.Ticket {
	unitID := ""
	if m.UnitID != nil {
		unitID = *m.UnitID
	}
	msgs := make([]*entity.TicketMessage, len(m.Messages))
	for i, msg := range m.Messages {
		msg := msg
		msgs[i] = msg.ToEntity()
	}
	return &entity.Ticket{
		ID:          m.ID,
		SchoolID:    m.SchoolID,
		UnitID:      unitID,
		RequesterID: m.RequesterID,
		Category:    entity.TicketCategory(m.Category),
		Status:      entity.TicketStatus(m.Status),
		Subject:     m.Subject,
		ResolvedAt:  m.ResolvedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Messages:    msgs,
	}
}

type TicketMessageModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TicketID  string    `gorm:"type:uuid;not null;index"`
	SenderID  string    `gorm:"type:uuid;not null"`
	Body      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (TicketMessageModel) TableName() string { return "ticket_messages" }

func (m *TicketMessageModel) ToEntity() *entity.TicketMessage {
	return &entity.TicketMessage{
		ID:        m.ID,
		TicketID:  m.TicketID,
		SenderID:  m.SenderID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
}
