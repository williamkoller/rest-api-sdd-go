package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domainrepo.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *entity.Ticket, firstMessage string) error {
	var unitID *string
	if ticket.UnitID != "" {
		unitID = &ticket.UnitID
	}
	ticketID := uuid.New().String()
	m := &model.TicketModel{
		ID:          ticketID,
		SchoolID:    ticket.SchoolID,
		UnitID:      unitID,
		RequesterID: ticket.RequesterID,
		Category:    string(ticket.Category),
		Status:      string(entity.TicketOpen),
		Subject:     ticket.Subject,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		msg := &model.TicketMessageModel{
			ID:       uuid.New().String(),
			TicketID: ticketID,
			SenderID: ticket.RequesterID,
			Body:     firstMessage,
		}
		return tx.Create(msg).Error
	})
	if err != nil {
		return fmt.Errorf("ticket repository: create: %w", err)
	}
	ticket.ID = ticketID
	return nil
}

func (r *ticketRepository) FindAll(ctx context.Context, schoolID string, filters domainrepo.TicketFilters, page, perPage int) ([]*entity.Ticket, int64, error) {
	var ms []model.TicketModel
	var total int64
	q := r.db.WithContext(ctx).Model(&model.TicketModel{}).Where("school_id = ?", schoolID)
	if filters.Status != "" {
		q = q.Where("status = ?", filters.Status)
	}
	if filters.Category != "" {
		q = q.Where("category = ?", filters.Category)
	}
	if filters.RequesterID != "" {
		q = q.Where("requester_id = ?", filters.RequesterID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("ticket repository: count: %w", err)
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("ticket repository: find all: %w", err)
	}
	result := make([]*entity.Ticket, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, total, nil
}

func (r *ticketRepository) FindByID(ctx context.Context, id string) (*entity.Ticket, error) {
	var m model.TicketModel
	if err := r.db.WithContext(ctx).Preload("Messages").First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("ticket repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *ticketRepository) AddMessage(ctx context.Context, msg *entity.TicketMessage) error {
	m := &model.TicketMessageModel{
		ID:       uuid.New().String(),
		TicketID: msg.TicketID,
		SenderID: msg.SenderID,
		Body:     msg.Body,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("ticket repository: add message: %w", err)
	}
	msg.ID = m.ID
	return nil
}

func (r *ticketRepository) UpdateStatus(ctx context.Context, id string, status entity.TicketStatus, resolvedAt *time.Time) error {
	updates := map[string]interface{}{"status": string(status)}
	if resolvedAt != nil {
		updates["resolved_at"] = resolvedAt
	}
	if err := r.db.WithContext(ctx).Model(&model.TicketModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("ticket repository: update status: %w", err)
	}
	return nil
}

func (r *ticketRepository) GetReport(ctx context.Context, schoolID string, from, to *time.Time) ([]*entity.TicketReport, error) {
	type reportRow struct {
		Category           string
		Count              int
		AvgResolutionHours float64
	}
	var rows []reportRow
	q := r.db.WithContext(ctx).Table("tickets").
		Select("category, COUNT(*) as count, AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))/3600) as avg_resolution_hours").
		Where("school_id = ?", schoolID)
	if from != nil {
		q = q.Where("created_at >= ?", from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", to)
	}
	if err := q.Group("category").Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("ticket repository: get report: %w", err)
	}
	result := make([]*entity.TicketReport, len(rows))
	for i, row := range rows {
		result[i] = &entity.TicketReport{
			Category:           row.Category,
			Count:              row.Count,
			AvgResolutionHours: row.AvgResolutionHours,
		}
	}
	return result, nil
}
