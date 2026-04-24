package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ReenrollmentCampaignResponse struct {
	ID           string    `json:"id"`
	SchoolID     string    `json:"schoolId"`
	UnitID       string    `json:"unitId"`
	AcademicYear int       `json:"academicYear"`
	Deadline     time.Time `json:"deadline"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ReenrollmentResponse struct {
	ID          string     `json:"id"`
	StudentID   string     `json:"studentId"`
	CampaignID  string     `json:"campaignId"`
	Status      string     `json:"status"`
	RespondedAt *time.Time `json:"respondedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type CampaignDashboardResponse struct {
	Total      int `json:"total"`
	Confirmed  int `json:"confirmed"`
	Declined   int `json:"declined"`
	NotStarted int `json:"notStarted"`
}

func MapReenrollmentCampaign(e *entity.ReenrollmentCampaign) *ReenrollmentCampaignResponse {
	return &ReenrollmentCampaignResponse{
		ID:           e.ID,
		SchoolID:     e.SchoolID,
		UnitID:       e.UnitID,
		AcademicYear: e.AcademicYear,
		Deadline:     e.Deadline,
		Status:       string(e.Status),
		CreatedAt:    e.CreatedAt,
	}
}

func MapReenrollment(e *entity.Reenrollment) *ReenrollmentResponse {
	return &ReenrollmentResponse{
		ID:          e.ID,
		StudentID:   e.StudentID,
		CampaignID:  e.CampaignID,
		Status:      string(e.Status),
		RespondedAt: e.RespondedAt,
		CreatedAt:   e.CreatedAt,
	}
}

func MapCampaignDashboard(e *entity.CampaignDashboard) *CampaignDashboardResponse {
	return &CampaignDashboardResponse{
		Total:      e.Total,
		Confirmed:  e.Confirmed,
		Declined:   e.Declined,
		NotStarted: e.NotStarted,
	}
}
