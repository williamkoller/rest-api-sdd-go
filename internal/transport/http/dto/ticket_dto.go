package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type TicketMessageResponse struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticketId"`
	SenderID  string    `json:"senderId"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type TicketResponse struct {
	ID          string                   `json:"id"`
	SchoolID    string                   `json:"schoolId"`
	UnitID      string                   `json:"unitId"`
	RequesterID string                   `json:"requesterId"`
	Category    string                   `json:"category"`
	Status      string                   `json:"status"`
	Subject     string                   `json:"subject"`
	ResolvedAt  *time.Time               `json:"resolvedAt"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
	Messages    []*TicketMessageResponse `json:"messages"`
}

type TicketReportResponse struct {
	Category           string  `json:"category"`
	Count              int     `json:"count"`
	AvgResolutionHours float64 `json:"avgResolutionHours"`
}

func MapTicketMessage(e *entity.TicketMessage) *TicketMessageResponse {
	return &TicketMessageResponse{
		ID:        e.ID,
		TicketID:  e.TicketID,
		SenderID:  e.SenderID,
		Body:      e.Body,
		CreatedAt: e.CreatedAt,
	}
}

func MapTicket(e *entity.Ticket) *TicketResponse {
	msgs := make([]*TicketMessageResponse, len(e.Messages))
	for i, m := range e.Messages {
		msgs[i] = MapTicketMessage(m)
	}
	return &TicketResponse{
		ID:          e.ID,
		SchoolID:    e.SchoolID,
		UnitID:      e.UnitID,
		RequesterID: e.RequesterID,
		Category:    string(e.Category),
		Status:      string(e.Status),
		Subject:     e.Subject,
		ResolvedAt:  e.ResolvedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		Messages:    msgs,
	}
}

func MapTickets(list []*entity.Ticket) []*TicketResponse {
	res := make([]*TicketResponse, len(list))
	for i, e := range list {
		res[i] = MapTicket(e)
	}
	return res
}

func MapTicketReport(list []*entity.TicketReport) []*TicketReportResponse {
	res := make([]*TicketReportResponse, len(list))
	for i, e := range list {
		res[i] = &TicketReportResponse{
			Category:           e.Category,
			Count:              e.Count,
			AvgResolutionHours: e.AvgResolutionHours,
		}
	}
	return res
}
