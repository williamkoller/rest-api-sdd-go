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

type waitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository(db *gorm.DB) domainrepo.WaitlistRepository {
	return &waitlistRepository{db: db}
}

func (r *waitlistRepository) Create(ctx context.Context, entry *entity.WaitlistEntry) error {
	var nextPos int
	r.db.WithContext(ctx).
		Table("waitlist_entries").
		Select("COALESCE(MAX(position), 0) + 1").
		Where("unit_id = ? AND grade_level = ? AND academic_year = ?", entry.UnitID, entry.GradeLevel, entry.AcademicYear).
		Scan(&nextPos)

	var referralID *string
	if entry.ReferralID != "" {
		referralID = &entry.ReferralID
	}

	m := &model.WaitlistEntryModel{
		ID:            uuid.New().String(),
		UnitID:        entry.UnitID,
		GuardianName:  entry.GuardianName,
		GuardianEmail: entry.GuardianEmail,
		StudentName:   entry.StudentName,
		GradeLevel:    entry.GradeLevel,
		AcademicYear:  entry.AcademicYear,
		Position:      nextPos,
		Status:        string(entity.WaitlistWaiting),
		ReferralID:    referralID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("waitlist repository: create: %w", err)
	}
	*entry = *m.ToEntity()
	return nil
}

func (r *waitlistRepository) FindByUnitID(ctx context.Context, unitID string, filters domainrepo.WaitlistFilters, page, perPage int) ([]*entity.WaitlistEntry, int64, error) {
	var ms []model.WaitlistEntryModel
	var total int64
	q := r.db.WithContext(ctx).Model(&model.WaitlistEntryModel{}).Where("unit_id = ?", unitID)
	if filters.Status != "" {
		q = q.Where("status = ?", filters.Status)
	}
	if filters.GradeLevel != "" {
		q = q.Where("grade_level = ?", filters.GradeLevel)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("waitlist repository: count: %w", err)
	}
	offset := (page - 1) * perPage
	if err := q.Order("position ASC").Offset(offset).Limit(perPage).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("waitlist repository: find by unit id: %w", err)
	}
	result := make([]*entity.WaitlistEntry, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, total, nil
}

func (r *waitlistRepository) FindByID(ctx context.Context, id string) (*entity.WaitlistEntry, error) {
	var m model.WaitlistEntryModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("waitlist repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *waitlistRepository) UpdateStatus(ctx context.Context, id string, status entity.WaitlistStatus) error {
	if err := r.db.WithContext(ctx).Model(&model.WaitlistEntryModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("waitlist repository: update status: %w", err)
	}
	return nil
}
