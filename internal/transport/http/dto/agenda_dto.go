package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AgendaItemResponse struct {
	ID          string     `json:"id"`
	ClassID     string     `json:"classId"`
	CreatedBy   string     `json:"createdBy"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func MapAgendaItem(e *entity.AgendaItem) *AgendaItemResponse {
	return &AgendaItemResponse{
		ID:          e.ID,
		ClassID:     e.ClassID,
		CreatedBy:   e.CreatedBy,
		Type:        string(e.Type),
		Title:       e.Title,
		Description: e.Description,
		DueDate:     e.DueDate,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func MapAgendaItems(list []*entity.AgendaItem) []*AgendaItemResponse {
	res := make([]*AgendaItemResponse, len(list))
	for i, e := range list {
		res[i] = MapAgendaItem(e)
	}
	return res
}
