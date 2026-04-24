package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type agendaRepository struct {
	db *gorm.DB
}

func NewAgendaRepository(db *gorm.DB) domainrepo.AgendaRepository {
	return &agendaRepository{db: db}
}

func (r *agendaRepository) Create(ctx context.Context, item *entity.AgendaItem) error {
	m := &model.AgendaItemModel{
		ID:          uuid.New().String(),
		ClassID:     item.ClassID,
		CreatedBy:   item.CreatedBy,
		Type:        string(item.Type),
		Title:       item.Title,
		Description: item.Description,
		DueDate:     item.DueDate,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("agenda repository: create: %w", err)
	}
	*item = *m.ToEntity()
	return nil
}

func (r *agendaRepository) FindByClassID(ctx context.Context, classID, typeFilter string, from, to *time.Time) ([]*entity.AgendaItem, error) {
	var ms []model.AgendaItemModel
	q := r.db.WithContext(ctx).Where("class_id = ?", classID)
	if typeFilter != "" {
		q = q.Where("type = ?", typeFilter)
	}
	if from != nil {
		q = q.Where("due_date >= ?", from)
	}
	if to != nil {
		q = q.Where("due_date <= ?", to)
	}
	if err := q.Order("due_date ASC").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("agenda repository: find by class id: %w", err)
	}
	result := make([]*entity.AgendaItem, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *agendaRepository) FindByID(ctx context.Context, id string) (*entity.AgendaItem, error) {
	var m model.AgendaItemModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("agenda repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *agendaRepository) Update(ctx context.Context, item *entity.AgendaItem) error {
	if err := r.db.WithContext(ctx).Model(&model.AgendaItemModel{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"type":        string(item.Type),
		"title":       item.Title,
		"description": item.Description,
		"due_date":    item.DueDate,
	}).Error; err != nil {
		return fmt.Errorf("agenda repository: update: %w", err)
	}
	return nil
}

func (r *agendaRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.AgendaItemModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("agenda repository: delete: %w", err)
	}
	return nil
}
