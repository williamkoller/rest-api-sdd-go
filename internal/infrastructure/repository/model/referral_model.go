package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReferralModel struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID      string    `gorm:"type:uuid;not null;index"`
	ReferrerID    string    `gorm:"type:uuid;not null;index"`
	ReferralCode  string    `gorm:"not null;uniqueIndex;size:8"`
	ReferredEmail string    `gorm:"size:200"`
	Status        string    `gorm:"not null;default:'pending'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (ReferralModel) TableName() string { return "referrals" }

func (m *ReferralModel) ToEntity() *entity.Referral {
	return &entity.Referral{
		ID:            m.ID,
		SchoolID:      m.SchoolID,
		ReferrerID:    m.ReferrerID,
		ReferralCode:  m.ReferralCode,
		ReferredEmail: m.ReferredEmail,
		Status:        entity.ReferralStatus(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
