package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ClassModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UnitID       string    `gorm:"type:uuid;not null;index"`
	ClassroomID  *string   `gorm:"type:uuid"`
	Name         string    `gorm:"not null;size:100"`
	GradeLevel   string    `gorm:"not null;size:20"`
	Shift        string    `gorm:"not null"`
	AcademicYear int       `gorm:"not null"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (ClassModel) TableName() string { return "classes" }

func (m *ClassModel) ToEntity() *entity.Class {
	classroomID := ""
	if m.ClassroomID != nil {
		classroomID = *m.ClassroomID
	}
	return &entity.Class{
		ID:           m.ID,
		UnitID:       m.UnitID,
		ClassroomID:  classroomID,
		Name:         m.Name,
		GradeLevel:   m.GradeLevel,
		Shift:        entity.Shift(m.Shift),
		AcademicYear: m.AcademicYear,
		Active:       m.Active,
		CreatedAt:    m.CreatedAt,
	}
}

func ClassFromEntity(c *entity.Class) *ClassModel {
	var classroomID *string
	if c.ClassroomID != "" {
		classroomID = &c.ClassroomID
	}
	return &ClassModel{
		ID:           c.ID,
		UnitID:       c.UnitID,
		ClassroomID:  classroomID,
		Name:         c.Name,
		GradeLevel:   c.GradeLevel,
		Shift:        string(c.Shift),
		AcademicYear: c.AcademicYear,
		Active:       c.Active,
		CreatedAt:    c.CreatedAt,
	}
}
