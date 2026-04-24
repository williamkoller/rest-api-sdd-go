package usecase

import (
	"context"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type ReferralUseCase struct {
	repo    repository.ReferralRepository
	baseURL string
}

func NewReferralUseCase(repo repository.ReferralRepository, baseURL string) *ReferralUseCase {
	return &ReferralUseCase{repo: repo, baseURL: baseURL}
}

func (uc *ReferralUseCase) GetMyLink(ctx context.Context, guardianID, schoolID string) (map[string]string, error) {
	referral, err := uc.repo.FindOrCreateByReferrer(ctx, schoolID, guardianID)
	if err != nil {
		return nil, fmt.Errorf("referral usecase: get link: %w", err)
	}
	return map[string]string{
		"referral_code": referral.ReferralCode,
		"link":          uc.baseURL + "/register?ref=" + referral.ReferralCode,
		"status":        string(referral.Status),
	}, nil
}

func (uc *ReferralUseCase) GetReport(ctx context.Context, schoolID, status, role string, page, perPage int) ([]*entity.Referral, int64, error) {
	if role != string(entity.RoleSchoolAdmin) && role != string(entity.RoleSuperAdmin) {
		return nil, 0, ErrForbidden
	}
	referrals, total, err := uc.repo.FindBySchoolID(ctx, schoolID, status, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("referral usecase: get report: %w", err)
	}
	return referrals, total, nil
}
