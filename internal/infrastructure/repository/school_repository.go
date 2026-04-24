package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type schoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) domainrepo.SchoolRepository {
	return &schoolRepository{db: db}
}

func (r *schoolRepository) Create(ctx context.Context, school *entity.School) error {
	m := model.SchoolFromEntity(school)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("school repository: create: %w", err)
	}
	*school = *m.ToEntity()
	return nil
}

func (r *schoolRepository) FindByID(ctx context.Context, id string) (*entity.School, error) {
	var m model.SchoolModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("school repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *schoolRepository) FindAll(ctx context.Context, page, perPage int, active *bool) ([]*entity.School, int64, error) {
	var ms []model.SchoolModel
	var total int64

	q := r.db.WithContext(ctx).Model(&model.SchoolModel{})
	if active != nil {
		q = q.Where("active = ?", *active)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("school repository: count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := q.Offset(offset).Limit(perPage).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("school repository: find all: %w", err)
	}

	schools := make([]*entity.School, len(ms))
	for i, m := range ms {
		m := m
		schools[i] = m.ToEntity()
	}
	return schools, total, nil
}

func (r *schoolRepository) Update(ctx context.Context, school *entity.School) error {
	m := model.SchoolFromEntity(school)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("school repository: update: %w", err)
	}
	return nil
}
