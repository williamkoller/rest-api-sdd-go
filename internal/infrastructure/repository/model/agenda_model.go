package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AgendaItemModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ClassID     string `gorm:"type:uuid;not null;index"`
	CreatedBy   string `gorm:"type:uuid;not null"`
	Type        string `gorm:"not null;size:20"`
	Title       string `gorm:"not null;size:300"`
	Description string `gorm:"type:text"`
	DueDate     *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (AgendaItemModel) TableName() string { return "agenda_items" }

func (m *AgendaItemModel) ToEntity() *entity.AgendaItem {
	return &entity.AgendaItem{
		ID:          m.ID,
		ClassID:     m.ClassID,
		CreatedBy:   m.CreatedBy,
		Type:        entity.AgendaItemType(m.Type),
		Title:       m.Title,
		Description: m.Description,
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
