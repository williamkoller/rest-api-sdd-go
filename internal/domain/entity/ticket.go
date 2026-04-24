package entity

import "time"

type TicketCategory string

const (
	TicketGeneral        TicketCategory = "general"
	TicketFinancial      TicketCategory = "financial"
	TicketAcademic       TicketCategory = "academic"
	TicketAdministrative TicketCategory = "administrative"
)

type TicketStatus string

const (
	TicketOpen       TicketStatus = "open"
	TicketInProgress TicketStatus = "in_progress"
	TicketResolved   TicketStatus = "resolved"
	TicketClosed     TicketStatus = "closed"
)

type Ticket struct {
	ID          string
	SchoolID    string
	UnitID      string
	RequesterID string
	Category    TicketCategory
	Status      TicketStatus
	Subject     string
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Messages    []*TicketMessage
}

type TicketMessage struct {
	ID        string
	TicketID  string
	SenderID  string
	Body      string
	CreatedAt time.Time
}

type TicketReport struct {
	Category           string
	Count              int
	AvgResolutionHours float64
}
