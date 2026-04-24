package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type UnitModel struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID  string `gorm:"type:uuid;not null;index"`
	Name      string `gorm:"not null"`
	Address   string `gorm:"not null"`
	City      string `gorm:"not null"`
	State     string `gorm:"not null;size:2"`
	ZipCode   string `gorm:"not null"`
	Phone     string
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UnitModel) TableName() string { return "units" }

func (m *UnitModel) ToEntity() *entity.Unit {
	return &entity.Unit{
		ID:        m.ID,
		SchoolID:  m.SchoolID,
		Name:      m.Name,
		Address:   m.Address,
		City:      m.City,
		State:     m.State,
		ZipCode:   m.ZipCode,
		Phone:     m.Phone,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func UnitFromEntity(u *entity.Unit) *UnitModel {
	return &UnitModel{
		ID:        u.ID,
		SchoolID:  u.SchoolID,
		Name:      u.Name,
		Address:   u.Address,
		City:      u.City,
		State:     u.State,
		ZipCode:   u.ZipCode,
		Phone:     u.Phone,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
