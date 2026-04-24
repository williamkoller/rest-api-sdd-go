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

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) domainrepo.UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) Create(ctx context.Context, unit *entity.Unit) error {
	m := model.UnitFromEntity(unit)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("unit repository: create: %w", err)
	}
	*unit = *m.ToEntity()
	return nil
}

func (r *unitRepository) FindByID(ctx context.Context, id string) (*entity.Unit, error) {
	var m model.UnitModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("unit repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *unitRepository) FindBySchoolID(ctx context.Context, schoolID string, active *bool) ([]*entity.Unit, error) {
	var ms []model.UnitModel
	q := r.db.WithContext(ctx).Where("school_id = ?", schoolID)
	if active != nil {
		q = q.Where("active = ?", *active)
	}
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("unit repository: find by school id: %w", err)
	}
	units := make([]*entity.Unit, len(ms))
	for i, m := range ms {
		m := m
		units[i] = m.ToEntity()
	}
	return units, nil
}

func (r *unitRepository) Update(ctx context.Context, unit *entity.Unit) error {
	m := model.UnitFromEntity(unit)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("unit repository: update: %w", err)
	}
	return nil
}
