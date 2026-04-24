package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReenrollmentRepository interface {
	CreateCampaign(ctx context.Context, campaign *entity.ReenrollmentCampaign) error
	FindCampaignByID(ctx context.Context, id string) (*entity.ReenrollmentCampaign, error)
	GetDashboard(ctx context.Context, campaignID string) (*entity.CampaignDashboard, error)
	UpsertResponse(ctx context.Context, campaignID, studentID string, status entity.ReenrollmentStatus) (*entity.Reenrollment, error)
	FindByCampaignAndStudent(ctx context.Context, campaignID, studentID string) (*entity.Reenrollment, error)
}
