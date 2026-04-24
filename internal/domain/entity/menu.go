package entity

import "time"

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
)

type MealType string

const (
	MealBreakfast MealType = "breakfast"
	MealLunch     MealType = "lunch"
	MealSnack     MealType = "snack"
	MealDinner    MealType = "dinner"
)

type Menu struct {
	ID        string
	UnitID    string
	WeekStart time.Time
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []*MenuItem
}

type MenuItem struct {
	ID          string
	MenuID      string
	DayOfWeek   DayOfWeek
	MealType    MealType
	Description string
}
