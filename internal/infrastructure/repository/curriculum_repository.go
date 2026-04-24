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

type curriculumRepository struct {
	db *gorm.DB
}

func NewCurriculumRepository(db *gorm.DB) domainrepo.CurriculumRepository {
	return &curriculumRepository{db: db}
}

func (r *curriculumRepository) BatchCreate(ctx context.Context, classID string, entries []*entity.CurriculumEntry) error {
	ms := make([]model.CurriculumEntryModel, len(entries))
	for i, e := range entries {
		ms[i] = model.CurriculumEntryModel{
			ID:        uuid.New().String(),
			ClassID:   classID,
			Subject:   e.Subject,
			TeacherID: e.TeacherID,
			DayOfWeek: string(e.DayOfWeek),
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
		}
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CurriculumEntryModel{}, "class_id = ?", classID).Error; err != nil {
			return err
		}
		if len(ms) > 0 {
			return tx.Create(&ms).Error
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("curriculum repository: batch create: %w", err)
	}
	return nil
}

func (r *curriculumRepository) FindByClassID(ctx context.Context, classID string) ([]*entity.CurriculumEntry, error) {
	var ms []model.CurriculumEntryModel
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).Order("day_of_week, start_time").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("curriculum repository: find by class id: %w", err)
	}
	result := make([]*entity.CurriculumEntry, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *curriculumRepository) FindByID(ctx context.Context, id string) (*entity.CurriculumEntry, error) {
	var m model.CurriculumEntryModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("curriculum repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *curriculumRepository) Update(ctx context.Context, entry *entity.CurriculumEntry) error {
	if err := r.db.WithContext(ctx).Model(&model.CurriculumEntryModel{}).Where("id = ?", entry.ID).Updates(map[string]interface{}{
		"subject":     entry.Subject,
		"teacher_id":  entry.TeacherID,
		"day_of_week": string(entry.DayOfWeek),
		"start_time":  entry.StartTime,
		"end_time":    entry.EndTime,
	}).Error; err != nil {
		return fmt.Errorf("curriculum repository: update: %w", err)
	}
	return nil
}

func (r *curriculumRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.CurriculumEntryModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("curriculum repository: delete: %w", err)
	}
	return nil
}
