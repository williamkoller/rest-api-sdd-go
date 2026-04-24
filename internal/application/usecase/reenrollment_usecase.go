package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var (
	ErrCampaignNotOpen = errors.New("campaign is not open")
	ErrDeadlinePassed  = errors.New("reenrollment deadline has passed")
	ErrOutstandingDebt = errors.New("student has outstanding debt")
)

type ReenrollmentUseCase struct {
	reenrollRepo repository.ReenrollmentRepository
	invoiceRepo  repository.InvoiceRepository
}

func NewReenrollmentUseCase(reenrollRepo repository.ReenrollmentRepository, invoiceRepo repository.InvoiceRepository) *ReenrollmentUseCase {
	return &ReenrollmentUseCase{reenrollRepo: reenrollRepo, invoiceRepo: invoiceRepo}
}

func (uc *ReenrollmentUseCase) CreateCampaign(ctx context.Context, schoolID, unitID string, academicYear int, deadline time.Time) (*entity.ReenrollmentCampaign, error) {
	campaign := &entity.ReenrollmentCampaign{
		SchoolID:     schoolID,
		UnitID:       unitID,
		AcademicYear: academicYear,
		Deadline:     deadline,
		Status:       entity.CampaignOpen,
	}
	if err := uc.reenrollRepo.CreateCampaign(ctx, campaign); err != nil {
		return nil, fmt.Errorf("reenrollment usecase: create campaign: %w", err)
	}
	return campaign, nil
}

func (uc *ReenrollmentUseCase) GetDashboard(ctx context.Context, campaignID string) (*entity.CampaignDashboard, error) {
	campaign, err := uc.reenrollRepo.FindCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("reenrollment usecase: get dashboard find campaign: %w", err)
	}
	if campaign == nil {
		return nil, ErrNotFound
	}
	dash, err := uc.reenrollRepo.GetDashboard(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("reenrollment usecase: get dashboard: %w", err)
	}
	return dash, nil
}

func (uc *ReenrollmentUseCase) Respond(ctx context.Context, campaignID, studentID string, status entity.ReenrollmentStatus) (*entity.Reenrollment, error) {
	campaign, err := uc.reenrollRepo.FindCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("reenrollment usecase: respond find campaign: %w", err)
	}
	if campaign == nil {
		return nil, ErrNotFound
	}
	if campaign.Status != entity.CampaignOpen {
		return nil, ErrCampaignNotOpen
	}
	if time.Now().After(campaign.Deadline) {
		return nil, ErrDeadlinePassed
	}
	if status == entity.ReenrollmentConfirmed {
		overdue, err := uc.invoiceRepo.HasOverdueInvoices(ctx, studentID)
		if err != nil {
			return nil, fmt.Errorf("reenrollment usecase: check overdue: %w", err)
		}
		if overdue {
			return nil, ErrOutstandingDebt
		}
	}
	record, err := uc.reenrollRepo.UpsertResponse(ctx, campaignID, studentID, status)
	if err != nil {
		return nil, fmt.Errorf("reenrollment usecase: respond: %w", err)
	}
	return record, nil
}
