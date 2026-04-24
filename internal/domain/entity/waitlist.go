package entity

import "time"

type WaitlistStatus string

const (
	WaitlistWaiting   WaitlistStatus = "waiting"
	WaitlistOfferMade WaitlistStatus = "offer_made"
	WaitlistAccepted  WaitlistStatus = "accepted"
	WaitlistDeclined  WaitlistStatus = "declined"
	WaitlistExpired   WaitlistStatus = "expired"
)

type WaitlistEntry struct {
	ID            string
	UnitID        string
	GuardianName  string
	GuardianEmail string
	StudentName   string
	GradeLevel    string
	AcademicYear  int
	Position      int
	Status        WaitlistStatus
	ReferralID    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
