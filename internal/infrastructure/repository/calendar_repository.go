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

type calendarRepository struct {
	db *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) domainrepo.CalendarRepository {
	return &calendarRepository{db: db}
}

func (r *calendarRepository) Create(ctx context.Context, event *entity.CalendarEvent) error {
	var unitID *string
	if event.UnitID != "" {
		unitID = &event.UnitID
	}
	m := &model.CalendarEventModel{
		ID:          uuid.New().String(),
		SchoolID:    event.SchoolID,
		UnitID:      unitID,
		Title:       event.Title,
		Description: event.Description,
		Type:        string(event.Type),
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		CreatedBy:   event.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("calendar repository: create: %w", err)
	}
	*event = *m.ToEntity()
	return nil
}

func (r *calendarRepository) FindBySchoolID(ctx context.Context, schoolID string, unitID *string, from, to *time.Time) ([]*entity.CalendarEvent, error) {
	var ms []model.CalendarEventModel
	q := r.db.WithContext(ctx).Where("school_id = ?", schoolID)
	if unitID != nil {
		q = q.Where("unit_id IS NULL OR unit_id = ?", *unitID)
	} else {
		q = q.Where("unit_id IS NULL")
	}
	if from != nil {
		q = q.Where("end_date >= ?", from)
	}
	if to != nil {
		q = q.Where("start_date <= ?", to)
	}
	if err := q.Order("start_date ASC").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("calendar repository: find by school id: %w", err)
	}
	result := make([]*entity.CalendarEvent, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *calendarRepository) FindByID(ctx context.Context, id string) (*entity.CalendarEvent, error) {
	var m model.CalendarEventModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("calendar repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *calendarRepository) Update(ctx context.Context, event *entity.CalendarEvent) error {
	var unitID *string
	if event.UnitID != "" {
		unitID = &event.UnitID
	}
	if err := r.db.WithContext(ctx).Model(&model.CalendarEventModel{}).Where("id = ?", event.ID).Updates(map[string]interface{}{
		"unit_id":     unitID,
		"title":       event.Title,
		"description": event.Description,
		"type":        string(event.Type),
		"start_date":  event.StartDate,
		"end_date":    event.EndDate,
	}).Error; err != nil {
		return fmt.Errorf("calendar repository: update: %w", err)
	}
	return nil
}

func (r *calendarRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.CalendarEventModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("calendar repository: delete: %w", err)
	}
	return nil
}
