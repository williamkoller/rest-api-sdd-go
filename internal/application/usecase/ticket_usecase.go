package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var staffRoles = map[string]bool{
	string(entity.RoleSchoolAdmin): true,
	string(entity.RoleUnitStaff):   true,
	string(entity.RoleSuperAdmin):  true,
}

type TicketUseCase struct {
	repo repository.TicketRepository
}

func NewTicketUseCase(repo repository.TicketRepository) *TicketUseCase {
	return &TicketUseCase{repo: repo}
}

func (uc *TicketUseCase) Create(ctx context.Context, ticket *entity.Ticket, firstMessage string) (*entity.Ticket, error) {
	if err := uc.repo.Create(ctx, ticket, firstMessage); err != nil {
		return nil, fmt.Errorf("ticket usecase: create: %w", err)
	}
	return ticket, nil
}

func (uc *TicketUseCase) List(ctx context.Context, schoolID, requesterID, role string, filters repository.TicketFilters, page, perPage int) ([]*entity.Ticket, int64, error) {
	if !staffRoles[role] {
		filters.RequesterID = requesterID
	}
	tickets, total, err := uc.repo.FindAll(ctx, schoolID, filters, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("ticket usecase: list: %w", err)
	}
	return tickets, total, nil
}

func (uc *TicketUseCase) GetByID(ctx context.Context, id, requesterID, role string) (*entity.Ticket, error) {
	ticket, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ticket usecase: get by id: %w", err)
	}
	if ticket == nil {
		return nil, ErrNotFound
	}
	if !staffRoles[role] && ticket.RequesterID != requesterID {
		return nil, ErrForbidden
	}
	return ticket, nil
}

func (uc *TicketUseCase) Reply(ctx context.Context, ticketID, body, senderID, role string) (*entity.TicketMessage, error) {
	ticket, err := uc.repo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket usecase: reply find: %w", err)
	}
	if ticket == nil {
		return nil, ErrNotFound
	}
	msg := &entity.TicketMessage{
		TicketID: ticketID,
		SenderID: senderID,
		Body:     body,
	}
	if err := uc.repo.AddMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("ticket usecase: reply add message: %w", err)
	}
	if staffRoles[role] && ticket.Status == entity.TicketOpen {
		if err := uc.repo.UpdateStatus(ctx, ticketID, entity.TicketInProgress, nil); err != nil {
			return nil, fmt.Errorf("ticket usecase: reply update status: %w", err)
		}
	}
	return msg, nil
}

func (uc *TicketUseCase) UpdateStatus(ctx context.Context, id string, status entity.TicketStatus, role string) error {
	if !staffRoles[role] {
		return ErrForbidden
	}
	ticket, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ticket usecase: update status find: %w", err)
	}
	if ticket == nil {
		return ErrNotFound
	}
	var resolvedAt *time.Time
	if status == entity.TicketResolved {
		now := time.Now()
		resolvedAt = &now
	}
	if err := uc.repo.UpdateStatus(ctx, id, status, resolvedAt); err != nil {
		return fmt.Errorf("ticket usecase: update status: %w", err)
	}
	return nil
}

func (uc *TicketUseCase) GetReport(ctx context.Context, schoolID string, from, to *time.Time) ([]*entity.TicketReport, error) {
	report, err := uc.repo.GetReport(ctx, schoolID, from, to)
	if err != nil {
		return nil, fmt.Errorf("ticket usecase: get report: %w", err)
	}
	return report, nil
}
