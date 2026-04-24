package entity

import "time"

type Shift string

const (
	ShiftMorning   Shift = "morning"
	ShiftAfternoon Shift = "afternoon"
	ShiftFull      Shift = "full"
)

type Class struct {
	ID           string
	UnitID       string
	ClassroomID  string
	Name         string
	GradeLevel   string
	Shift        Shift
	AcademicYear int
	Active       bool
	CreatedAt    time.Time
}
