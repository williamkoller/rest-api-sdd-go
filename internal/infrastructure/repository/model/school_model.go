package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type SchoolModel struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string `gorm:"not null"`
	CNPJ      string `gorm:"not null;uniqueIndex"`
	Email     string `gorm:"not null"`
	Phone     string
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (SchoolModel) TableName() string { return "schools" }

func (m *SchoolModel) ToEntity() *entity.School {
	return &entity.School{
		ID:        m.ID,
		Name:      m.Name,
		CNPJ:      m.CNPJ,
		Email:     m.Email,
		Phone:     m.Phone,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func SchoolFromEntity(s *entity.School) *SchoolModel {
	return &SchoolModel{
		ID:        s.ID,
		Name:      s.Name,
		CNPJ:      s.CNPJ,
		Email:     s.Email,
		Phone:     s.Phone,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
