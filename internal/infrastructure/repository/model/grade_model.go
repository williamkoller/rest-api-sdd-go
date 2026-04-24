package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type GradeModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EnrollmentID string    `gorm:"type:uuid;not null;index"`
	Subject      string    `gorm:"not null;size:100"`
	Period       string    `gorm:"not null;size:20"`
	Value        float64   `gorm:"type:decimal(5,2);not null"`
	RecordedBy   string    `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (GradeModel) TableName() string { return "grades" }

func (m *GradeModel) ToEntity() *entity.Grade {
	return &entity.Grade{
		ID:           m.ID,
		EnrollmentID: m.EnrollmentID,
		Subject:      m.Subject,
		Period:       m.Period,
		Value:        m.Value,
		RecordedBy:   m.RecordedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
