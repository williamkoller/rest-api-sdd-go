package entity

import "time"

type CalendarEventType string

const (
	CalendarHoliday    CalendarEventType = "holiday"
	CalendarExamPeriod CalendarEventType = "exam_period"
	CalendarTypeEvent  CalendarEventType = "event"
	CalendarRecess     CalendarEventType = "recess"
)

type CalendarEvent struct {
	ID          string
	SchoolID    string
	UnitID      string
	Title       string
	Description string
	Type        CalendarEventType
	StartDate   time.Time
	EndDate     time.Time
	CreatedBy   string
	CreatedAt   time.Time
}
