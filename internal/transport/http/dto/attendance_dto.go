package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AttendanceRecordResponse struct {
	ID           string    `json:"id"`
	EnrollmentID string    `json:"enrollmentId"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	Note         string    `json:"note"`
	RecordedBy   string    `json:"recordedBy"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AttendanceSummaryResponse struct {
	Records []*AttendanceRecordResponse `json:"records"`
	Summary AttendanceSummaryStats      `json:"summary"`
}

type AttendanceSummaryStats struct {
	Rate    float64 `json:"rate"`
	Total   int     `json:"total"`
	Present int     `json:"present"`
}

func MapAttendanceRecord(e *entity.AttendanceRecord) *AttendanceRecordResponse {
	return &AttendanceRecordResponse{
		ID:           e.ID,
		EnrollmentID: e.EnrollmentID,
		Date:         e.Date,
		Status:       string(e.Status),
		Note:         e.Note,
		RecordedBy:   e.RecordedBy,
		CreatedAt:    e.CreatedAt,
	}
}

func MapAttendanceSummary(s *entity.AttendanceSummary) *AttendanceSummaryResponse {
	records := make([]*AttendanceRecordResponse, len(s.Records))
	for i, r := range s.Records {
		records[i] = MapAttendanceRecord(r)
	}
	return &AttendanceSummaryResponse{
		Records: records,
		Summary: AttendanceSummaryStats{Rate: s.Rate, Total: s.Total, Present: s.Present},
	}
}
