package entity

import "time"

type AttendanceStatus string

const (
	AttendancePresent   AttendanceStatus = "present"
	AttendanceAbsent    AttendanceStatus = "absent"
	AttendanceJustified AttendanceStatus = "justified"
)

type AttendanceRecord struct {
	ID           string
	EnrollmentID string
	Date         time.Time
	Status       AttendanceStatus
	Note         string
	RecordedBy   string
	CreatedAt    time.Time
}

type AttendanceSummary struct {
	Records []*AttendanceRecord
	Rate    float64
	Total   int
	Present int
}
