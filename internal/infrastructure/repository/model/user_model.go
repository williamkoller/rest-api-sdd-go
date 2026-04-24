package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type UserModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID     *string   `gorm:"type:uuid;index"`
	Name         string    `gorm:"not null;size:200"`
	Email        string    `gorm:"not null;uniqueIndex;size:255"`
	PasswordHash string    `gorm:"not null;size:255"`
	Role         string    `gorm:"not null"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string { return "users" }

func (m *UserModel) ToEntity() *entity.User {
	schoolID := ""
	if m.SchoolID != nil {
		schoolID = *m.SchoolID
	}
	return &entity.User{
		ID:           m.ID,
		SchoolID:     schoolID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         entity.Role(m.Role),
		Active:       m.Active,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func UserFromEntity(u *entity.User) *UserModel {
	var schoolID *string
	if u.SchoolID != "" {
		schoolID = &u.SchoolID
	}
	return &UserModel{
		ID:           u.ID,
		SchoolID:     schoolID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		Active:       u.Active,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
