package entity

import "time"

type EnrollmentStatus string

const (
	EnrollmentStatusActive      EnrollmentStatus = "active"
	EnrollmentStatusTransferred EnrollmentStatus = "transferred"
	EnrollmentStatusUnenrolled  EnrollmentStatus = "unenrolled"
)

type Enrollment struct {
	ID           string
	StudentID    string
	ClassID      string
	AcademicYear int
	EnrolledAt   time.Time
	UnenrolledAt *time.Time
	Status       EnrollmentStatus
}
