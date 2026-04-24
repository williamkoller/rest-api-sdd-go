package entity

import "time"

type Grade struct {
	ID           string
	EnrollmentID string
	Subject      string
	Period       string
	Value        float64
	RecordedBy   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
