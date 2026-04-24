package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type CalendarEventModel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID    string    `gorm:"type:uuid;not null;index"`
	UnitID      *string   `gorm:"type:uuid"`
	Title       string    `gorm:"not null;size:300"`
	Description string    `gorm:"type:text"`
	Type        string    `gorm:"not null;size:20"`
	StartDate   time.Time `gorm:"type:date;not null"`
	EndDate     time.Time `gorm:"type:date;not null"`
	CreatedBy   string    `gorm:"type:uuid;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (CalendarEventModel) TableName() string { return "calendar_events" }

func (m *CalendarEventModel) ToEntity() *entity.CalendarEvent {
	unitID := ""
	if m.UnitID != nil {
		unitID = *m.UnitID
	}
	return &entity.CalendarEvent{
		ID:          m.ID,
		SchoolID:    m.SchoolID,
		UnitID:      unitID,
		Title:       m.Title,
		Description: m.Description,
		Type:        entity.CalendarEventType(m.Type),
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
	}
}
