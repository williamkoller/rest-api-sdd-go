package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type GradeResponse struct {
	ID           string    `json:"id"`
	EnrollmentID string    `json:"enrollmentId"`
	Subject      string    `json:"subject"`
	Period       string    `json:"period"`
	Value        float64   `json:"value"`
	RecordedBy   string    `json:"recordedBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func MapGrade(e *entity.Grade) *GradeResponse {
	return &GradeResponse{
		ID:           e.ID,
		EnrollmentID: e.EnrollmentID,
		Subject:      e.Subject,
		Period:       e.Period,
		Value:        e.Value,
		RecordedBy:   e.RecordedBy,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func MapGrades(list []*entity.Grade) []*GradeResponse {
	res := make([]*GradeResponse, len(list))
	for i, e := range list {
		res[i] = MapGrade(e)
	}
	return res
}
