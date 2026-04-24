package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type TicketFilters struct {
	Status      string
	Category    string
	RequesterID string
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *entity.Ticket, firstMessage string) error
	FindAll(ctx context.Context, schoolID string, filters TicketFilters, page, perPage int) ([]*entity.Ticket, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Ticket, error)
	AddMessage(ctx context.Context, msg *entity.TicketMessage) error
	UpdateStatus(ctx context.Context, id string, status entity.TicketStatus, resolvedAt *time.Time) error
	GetReport(ctx context.Context, schoolID string, from, to *time.Time) ([]*entity.TicketReport, error)
}
