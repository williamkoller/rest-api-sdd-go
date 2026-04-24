package entity

import "time"

type CampaignStatus string

const (
	CampaignOpen   CampaignStatus = "open"
	CampaignClosed CampaignStatus = "closed"
)

type ReenrollmentStatus string

const (
	ReenrollmentNotStarted ReenrollmentStatus = "not_started"
	ReenrollmentConfirmed  ReenrollmentStatus = "confirmed"
	ReenrollmentDeclined   ReenrollmentStatus = "declined"
)

type ReenrollmentCampaign struct {
	ID           string
	SchoolID     string
	UnitID       string
	AcademicYear int
	Deadline     time.Time
	Status       CampaignStatus
	CreatedAt    time.Time
}

type Reenrollment struct {
	ID          string
	StudentID   string
	CampaignID  string
	Status      ReenrollmentStatus
	RespondedAt *time.Time
	CreatedAt   time.Time
}

type CampaignDashboard struct {
	Total      int
	Confirmed  int
	Declined   int
	NotStarted int
}
