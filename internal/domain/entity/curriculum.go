package entity

type CurriculumEntry struct {
	ID        string
	ClassID   string
	Subject   string
	TeacherID string
	DayOfWeek DayOfWeek
	StartTime string
	EndTime   string
}
