package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type FeedPostModel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SchoolID    string    `gorm:"type:uuid;not null;index"`
	UnitID      *string   `gorm:"type:uuid"`
	AuthorID    string    `gorm:"type:uuid;not null"`
	Body        string    `gorm:"type:text;not null"`
	ImageURL    string    `gorm:"size:500"`
	PublishedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (FeedPostModel) TableName() string { return "feed_posts" }

func (m *FeedPostModel) ToEntity() *entity.FeedPost {
	unitID := ""
	if m.UnitID != nil {
		unitID = *m.UnitID
	}
	return &entity.FeedPost{
		ID:          m.ID,
		SchoolID:    m.SchoolID,
		UnitID:      unitID,
		AuthorID:    m.AuthorID,
		Body:        m.Body,
		ImageURL:    m.ImageURL,
		PublishedAt: m.PublishedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
