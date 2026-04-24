package model

import (
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type CurriculumEntryModel struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ClassID   string `gorm:"type:uuid;not null;index"`
	Subject   string `gorm:"not null;size:100"`
	TeacherID string `gorm:"type:uuid;not null"`
	DayOfWeek string `gorm:"not null;size:10"`
	StartTime string `gorm:"type:time;not null"`
	EndTime   string `gorm:"type:time;not null"`
}

func (CurriculumEntryModel) TableName() string { return "curriculum_entries" }

func (m *CurriculumEntryModel) ToEntity() *entity.CurriculumEntry {
	return &entity.CurriculumEntry{
		ID:        m.ID,
		ClassID:   m.ClassID,
		Subject:   m.Subject,
		TeacherID: m.TeacherID,
		DayOfWeek: entity.DayOfWeek(m.DayOfWeek),
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
	}
}
