package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AttendanceRecordModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EnrollmentID string    `gorm:"type:uuid;not null;index"`
	Date         time.Time `gorm:"not null"`
	Status       string    `gorm:"not null"`
	Note         string
	RecordedBy   string    `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (AttendanceRecordModel) TableName() string { return "attendance_records" }

func (m *AttendanceRecordModel) ToEntity() *entity.AttendanceRecord {
	return &entity.AttendanceRecord{
		ID:           m.ID,
		EnrollmentID: m.EnrollmentID,
		Date:         m.Date,
		Status:       entity.AttendanceStatus(m.Status),
		Note:         m.Note,
		RecordedBy:   m.RecordedBy,
		CreatedAt:    m.CreatedAt,
	}
}
