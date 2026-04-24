package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type EnrollmentModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	StudentID    string    `gorm:"type:uuid;not null;index"`
	ClassID      string    `gorm:"type:uuid;not null;index"`
	AcademicYear int       `gorm:"not null"`
	EnrolledAt   time.Time `gorm:"not null"`
	UnenrolledAt *time.Time
	Status       string `gorm:"not null;default:'active'"`
}

func (EnrollmentModel) TableName() string { return "enrollments" }

func (m *EnrollmentModel) ToEntity() *entity.Enrollment {
	return &entity.Enrollment{
		ID:           m.ID,
		StudentID:    m.StudentID,
		ClassID:      m.ClassID,
		AcademicYear: m.AcademicYear,
		EnrolledAt:   m.EnrolledAt,
		UnenrolledAt: m.UnenrolledAt,
		Status:       entity.EnrollmentStatus(m.Status),
	}
}

func EnrollmentFromEntity(e *entity.Enrollment) *EnrollmentModel {
	return &EnrollmentModel{
		ID:           e.ID,
		StudentID:    e.StudentID,
		ClassID:      e.ClassID,
		AcademicYear: e.AcademicYear,
		EnrolledAt:   e.EnrolledAt,
		UnenrolledAt: e.UnenrolledAt,
		Status:       string(e.Status),
	}
}
