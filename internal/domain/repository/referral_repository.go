package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReferralRepository interface {
	FindOrCreateByReferrer(ctx context.Context, schoolID, referrerID string) (*entity.Referral, error)
	FindByCode(ctx context.Context, code string) (*entity.Referral, error)
	UpdateStatus(ctx context.Context, id string, status entity.ReferralStatus) error
	FindBySchoolID(ctx context.Context, schoolID, status string, page, perPage int) ([]*entity.Referral, int64, error)
}
