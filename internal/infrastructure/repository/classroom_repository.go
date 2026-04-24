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

type classroomRepository struct {
	db *gorm.DB
}

func NewClassroomRepository(db *gorm.DB) domainrepo.ClassroomRepository {
	return &classroomRepository{db: db}
}

func (r *classroomRepository) Create(ctx context.Context, classroom *entity.Classroom) error {
	m := model.ClassroomFromEntity(classroom)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("classroom repository: create: %w", err)
	}
	*classroom = *m.ToEntity()
	return nil
}

func (r *classroomRepository) FindByID(ctx context.Context, id string) (*entity.Classroom, error) {
	var m model.ClassroomModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("classroom repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *classroomRepository) FindByUnitID(ctx context.Context, unitID string) ([]*entity.Classroom, error) {
	var ms []model.ClassroomModel
	if err := r.db.WithContext(ctx).Where("unit_id = ?", unitID).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("classroom repository: find by unit id: %w", err)
	}
	classrooms := make([]*entity.Classroom, len(ms))
	for i, m := range ms {
		m := m
		classrooms[i] = m.ToEntity()
	}
	return classrooms, nil
}

func (r *classroomRepository) Update(ctx context.Context, classroom *entity.Classroom) error {
	m := model.ClassroomFromEntity(classroom)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("classroom repository: update: %w", err)
	}
	return nil
}
