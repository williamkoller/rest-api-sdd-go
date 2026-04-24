package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const referralCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type referralRepository struct {
	db *gorm.DB
}

func NewReferralRepository(db *gorm.DB) domainrepo.ReferralRepository {
	return &referralRepository{db: db}
}

func generateCode() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = referralCodeChars[rand.Intn(len(referralCodeChars))]
	}
	return string(b)
}

func (r *referralRepository) FindOrCreateByReferrer(ctx context.Context, schoolID, referrerID string) (*entity.Referral, error) {
	var m model.ReferralModel
	err := r.db.WithContext(ctx).Where("school_id = ? AND referrer_id = ?", schoolID, referrerID).First(&m).Error
	if err == nil {
		return m.ToEntity(), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("referral repository: find by referrer: %w", err)
	}
	newModel := model.ReferralModel{
		ID:           uuid.New().String(),
		SchoolID:     schoolID,
		ReferrerID:   referrerID,
		ReferralCode: generateCode(),
		Status:       string(entity.ReferralPending),
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&newModel).Error; err != nil {
		return nil, fmt.Errorf("referral repository: create: %w", err)
	}
	return newModel.ToEntity(), nil
}

func (r *referralRepository) FindByCode(ctx context.Context, code string) (*entity.Referral, error) {
	var m model.ReferralModel
	if err := r.db.WithContext(ctx).Where("referral_code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("referral repository: find by code: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *referralRepository) UpdateStatus(ctx context.Context, id string, status entity.ReferralStatus) error {
	if err := r.db.WithContext(ctx).Model(&model.ReferralModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("referral repository: update status: %w", err)
	}
	return nil
}

func (r *referralRepository) FindBySchoolID(ctx context.Context, schoolID, status string, page, perPage int) ([]*entity.Referral, int64, error) {
	var ms []model.ReferralModel
	var total int64
	q := r.db.WithContext(ctx).Model(&model.ReferralModel{}).Where("school_id = ?", schoolID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("referral repository: count: %w", err)
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("referral repository: find by school: %w", err)
	}
	result := make([]*entity.Referral, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, total, nil
}
