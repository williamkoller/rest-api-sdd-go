package entity

import "time"

type ReferralStatus string

const (
	ReferralPending    ReferralStatus = "pending"
	ReferralRegistered ReferralStatus = "registered"
	ReferralEnrolled   ReferralStatus = "enrolled"
)

type Referral struct {
	ID            string
	SchoolID      string
	ReferrerID    string
	ReferralCode  string
	ReferredEmail string
	Status        ReferralStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
