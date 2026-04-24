package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) domainrepo.FeedRepository {
	return &feedRepository{db: db}
}

func (r *feedRepository) Create(ctx context.Context, post *entity.FeedPost) error {
	var unitID *string
	if post.UnitID != "" {
		unitID = &post.UnitID
	}
	m := &model.FeedPostModel{
		ID:       uuid.New().String(),
		SchoolID: post.SchoolID,
		UnitID:   unitID,
		AuthorID: post.AuthorID,
		Body:     post.Body,
		ImageURL: post.ImageURL,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("feed repository: create: %w", err)
	}
	*post = *m.ToEntity()
	return nil
}

func (r *feedRepository) FindBySchoolID(ctx context.Context, schoolID string, unitID *string, page, perPage int) ([]*entity.FeedPost, int64, error) {
	var ms []model.FeedPostModel
	var total int64
	q := r.db.WithContext(ctx).Model(&model.FeedPostModel{}).Where("school_id = ?", schoolID)
	if unitID != nil {
		q = q.Where("unit_id IS NULL OR unit_id = ?", *unitID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("feed repository: count: %w", err)
	}
	offset := (page - 1) * perPage
	if err := q.Order("published_at DESC").Offset(offset).Limit(perPage).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("feed repository: find by school id: %w", err)
	}
	result := make([]*entity.FeedPost, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, total, nil
}

func (r *feedRepository) FindByID(ctx context.Context, id string) (*entity.FeedPost, error) {
	var m model.FeedPostModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("feed repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *feedRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.FeedPostModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("feed repository: delete: %w", err)
	}
	return nil
}
