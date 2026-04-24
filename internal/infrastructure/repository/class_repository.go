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

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) domainrepo.ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(ctx context.Context, class *entity.Class) error {
	m := model.ClassFromEntity(class)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("class repository: create: %w", err)
	}
	*class = *m.ToEntity()
	return nil
}

func (r *classRepository) FindByID(ctx context.Context, id string) (*entity.Class, error) {
	var m model.ClassModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("class repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *classRepository) FindByUnitID(ctx context.Context, unitID string, year *int, active *bool) ([]*entity.Class, error) {
	var ms []model.ClassModel
	q := r.db.WithContext(ctx).Where("unit_id = ?", unitID)
	if year != nil {
		q = q.Where("academic_year = ?", *year)
	}
	if active != nil {
		q = q.Where("active = ?", *active)
	}
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("class repository: find by unit id: %w", err)
	}
	classes := make([]*entity.Class, len(ms))
	for i, m := range ms {
		m := m
		classes[i] = m.ToEntity()
	}
	return classes, nil
}

func (r *classRepository) Update(ctx context.Context, class *entity.Class) error {
	m := model.ClassFromEntity(class)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("class repository: update: %w", err)
	}
	return nil
}

func (r *classRepository) CountEnrollments(ctx context.Context, classID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.EnrollmentModel{}).
		Where("class_id = ? AND status = 'active'", classID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("class repository: count enrollments: %w", err)
	}
	return count, nil
}
