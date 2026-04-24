package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type CalendarEventResponse struct {
	ID          string    `json:"id"`
	SchoolID    string    `json:"schoolId"`
	UnitID      string    `json:"unitId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

func MapCalendarEvent(e *entity.CalendarEvent) *CalendarEventResponse {
	return &CalendarEventResponse{
		ID:          e.ID,
		SchoolID:    e.SchoolID,
		UnitID:      e.UnitID,
		Title:       e.Title,
		Description: e.Description,
		Type:        string(e.Type),
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		CreatedBy:   e.CreatedBy,
		CreatedAt:   e.CreatedAt,
	}
}

func MapCalendarEvents(list []*entity.CalendarEvent) []*CalendarEventResponse {
	res := make([]*CalendarEventResponse, len(list))
	for i, e := range list {
		res[i] = MapCalendarEvent(e)
	}
	return res
}
