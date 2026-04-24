package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type StudentModel struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID           string    `gorm:"type:uuid;not null;index"`
	Name               string    `gorm:"not null;size:200"`
	BirthDate          time.Time `gorm:"not null"`
	CPF                string    `gorm:"size:11"`
	RegistrationNumber string    `gorm:"not null;size:50"`
	Active             bool      `gorm:"not null;default:true"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (StudentModel) TableName() string { return "students" }

func (m *StudentModel) ToEntity() *entity.Student {
	return &entity.Student{
		ID:                 m.ID,
		SchoolID:           m.SchoolID,
		Name:               m.Name,
		BirthDate:          m.BirthDate,
		CPF:                m.CPF,
		RegistrationNumber: m.RegistrationNumber,
		Active:             m.Active,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func StudentFromEntity(s *entity.Student) *StudentModel {
	return &StudentModel{
		ID:                 s.ID,
		SchoolID:           s.SchoolID,
		Name:               s.Name,
		BirthDate:          s.BirthDate,
		CPF:                s.CPF,
		RegistrationNumber: s.RegistrationNumber,
		Active:             s.Active,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}
