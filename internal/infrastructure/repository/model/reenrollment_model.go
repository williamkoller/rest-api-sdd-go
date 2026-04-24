package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReenrollmentCampaignModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID     string    `gorm:"type:uuid;not null"`
	UnitID       *string   `gorm:"type:uuid"`
	AcademicYear int       `gorm:"not null"`
	Deadline     time.Time `gorm:"not null"`
	Status       string    `gorm:"not null;default:'open'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (ReenrollmentCampaignModel) TableName() string { return "reenrollment_campaigns" }

func (m *ReenrollmentCampaignModel) ToEntity() *entity.ReenrollmentCampaign {
	unitID := ""
	if m.UnitID != nil {
		unitID = *m.UnitID
	}
	return &entity.ReenrollmentCampaign{
		ID:           m.ID,
		SchoolID:     m.SchoolID,
		UnitID:       unitID,
		AcademicYear: m.AcademicYear,
		Deadline:     m.Deadline,
		Status:       entity.CampaignStatus(m.Status),
		CreatedAt:    m.CreatedAt,
	}
}

type ReenrollmentModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	StudentID   string `gorm:"type:uuid;not null"`
	CampaignID  string `gorm:"type:uuid;not null;index"`
	Status      string `gorm:"not null;default:'not_started'"`
	RespondedAt *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (ReenrollmentModel) TableName() string { return "reenrollments" }

func (m *ReenrollmentModel) ToEntity() *entity.Reenrollment {
	return &entity.Reenrollment{
		ID:          m.ID,
		StudentID:   m.StudentID,
		CampaignID:  m.CampaignID,
		Status:      entity.ReenrollmentStatus(m.Status),
		RespondedAt: m.RespondedAt,
		CreatedAt:   m.CreatedAt,
	}
}
