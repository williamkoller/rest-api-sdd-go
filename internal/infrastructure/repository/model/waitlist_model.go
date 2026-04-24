package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type WaitlistEntryModel struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UnitID        string    `gorm:"type:uuid;not null;index"`
	GuardianName  string    `gorm:"not null;size:200"`
	GuardianEmail string    `gorm:"not null;size:200"`
	StudentName   string    `gorm:"not null;size:200"`
	GradeLevel    string    `gorm:"not null;size:50"`
	AcademicYear  int       `gorm:"not null"`
	Position      int       `gorm:"not null"`
	Status        string    `gorm:"not null;default:'waiting'"`
	ReferralID    *string   `gorm:"type:uuid"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (WaitlistEntryModel) TableName() string { return "waitlist_entries" }

func (m *WaitlistEntryModel) ToEntity() *entity.WaitlistEntry {
	referralID := ""
	if m.ReferralID != nil {
		referralID = *m.ReferralID
	}
	return &entity.WaitlistEntry{
		ID:            m.ID,
		UnitID:        m.UnitID,
		GuardianName:  m.GuardianName,
		GuardianEmail: m.GuardianEmail,
		StudentName:   m.StudentName,
		GradeLevel:    m.GradeLevel,
		AcademicYear:  m.AcademicYear,
		Position:      m.Position,
		Status:        entity.WaitlistStatus(m.Status),
		ReferralID:    referralID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
