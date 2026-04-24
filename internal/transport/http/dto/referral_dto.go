package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReferralResponse struct {
	ID            string    `json:"id"`
	SchoolID      string    `json:"schoolId"`
	ReferrerID    string    `json:"referrerId"`
	ReferralCode  string    `json:"referralCode"`
	ReferredEmail string    `json:"referredEmail"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func MapReferral(e *entity.Referral) *ReferralResponse {
	return &ReferralResponse{
		ID:            e.ID,
		SchoolID:      e.SchoolID,
		ReferrerID:    e.ReferrerID,
		ReferralCode:  e.ReferralCode,
		ReferredEmail: e.ReferredEmail,
		Status:        string(e.Status),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func MapReferrals(list []*entity.Referral) []*ReferralResponse {
	res := make([]*ReferralResponse, len(list))
	for i, e := range list {
		res[i] = MapReferral(e)
	}
	return res
}
