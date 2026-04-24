package entity

import "time"

type AgendaItemType string

const (
	AgendaHomework AgendaItemType = "homework"
	AgendaEvent    AgendaItemType = "event"
	AgendaReminder AgendaItemType = "reminder"
)

type AgendaItem struct {
	ID          string
	ClassID     string
	CreatedBy   string
	Type        AgendaItemType
	Title       string
	Description string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
