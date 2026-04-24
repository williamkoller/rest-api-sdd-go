package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type MenuItemResponse struct {
	ID          string `json:"id"`
	MenuID      string `json:"menuId"`
	DayOfWeek   string `json:"dayOfWeek"`
	MealType    string `json:"mealType"`
	Description string `json:"description"`
}

type MenuResponse struct {
	ID        string              `json:"id"`
	UnitID    string              `json:"unitId"`
	WeekStart time.Time           `json:"weekStart"`
	CreatedBy string              `json:"createdBy"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
	Items     []*MenuItemResponse `json:"items"`
}

func MapMenuItem(e *entity.MenuItem) *MenuItemResponse {
	return &MenuItemResponse{
		ID:          e.ID,
		MenuID:      e.MenuID,
		DayOfWeek:   string(e.DayOfWeek),
		MealType:    string(e.MealType),
		Description: e.Description,
	}
}

func MapMenu(e *entity.Menu) *MenuResponse {
	items := make([]*MenuItemResponse, len(e.Items))
	for i, item := range e.Items {
		items[i] = MapMenuItem(item)
	}
	return &MenuResponse{
		ID:        e.ID,
		UnitID:    e.UnitID,
		WeekStart: e.WeekStart,
		CreatedBy: e.CreatedBy,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Items:     items,
	}
}
