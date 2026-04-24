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
	"gorm.io/gorm/clause"
)

type reenrollmentRepository struct {
	db *gorm.DB
}

func NewReenrollmentRepository(db *gorm.DB) domainrepo.ReenrollmentRepository {
	return &reenrollmentRepository{db: db}
}

func (r *reenrollmentRepository) CreateCampaign(ctx context.Context, campaign *entity.ReenrollmentCampaign) error {
	var unitID *string
	if campaign.UnitID != "" {
		unitID = &campaign.UnitID
	}
	m := &model.ReenrollmentCampaignModel{
		ID:           uuid.New().String(),
		SchoolID:     campaign.SchoolID,
		UnitID:       unitID,
		AcademicYear: campaign.AcademicYear,
		Deadline:     campaign.Deadline,
		Status:       string(campaign.Status),
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("reenrollment repository: create campaign: %w", err)
	}
	*campaign = *m.ToEntity()
	return nil
}

func (r *reenrollmentRepository) FindCampaignByID(ctx context.Context, id string) (*entity.ReenrollmentCampaign, error) {
	var m model.ReenrollmentCampaignModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("reenrollment repository: find campaign: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *reenrollmentRepository) GetDashboard(ctx context.Context, campaignID string) (*entity.CampaignDashboard, error) {
	type row struct {
		Status string
		Count  int
	}
	var rows []row
	err := r.db.WithContext(ctx).Table("reenrollments").
		Select("status, COUNT(*) as count").
		Where("campaign_id = ?", campaignID).
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("reenrollment repository: dashboard: %w", err)
	}

	dash := &entity.CampaignDashboard{}
	for _, row := range rows {
		switch row.Status {
		case "confirmed":
			dash.Confirmed = row.Count
		case "declined":
			dash.Declined = row.Count
		case "not_started":
			dash.NotStarted = row.Count
		}
		dash.Total += row.Count
	}
	return dash, nil
}

func (r *reenrollmentRepository) UpsertResponse(ctx context.Context, campaignID, studentID string, status entity.ReenrollmentStatus) (*entity.Reenrollment, error) {
	now := time.Now()
	m := model.ReenrollmentModel{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		CampaignID:  campaignID,
		Status:      string(status),
		RespondedAt: &now,
	}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "student_id"}, {Name: "campaign_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "responded_at"}),
		}).
		Create(&m).Error
	if err != nil {
		return nil, fmt.Errorf("reenrollment repository: upsert response: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *reenrollmentRepository) FindByCampaignAndStudent(ctx context.Context, campaignID, studentID string) (*entity.Reenrollment, error) {
	var m model.ReenrollmentModel
	err := r.db.WithContext(ctx).Where("campaign_id = ? AND student_id = ?", campaignID, studentID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("reenrollment repository: find by campaign and student: %w", err)
	}
	return m.ToEntity(), nil
}
