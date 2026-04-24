package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type CalendarRepository interface {
	Create(ctx context.Context, event *entity.CalendarEvent) error
	FindBySchoolID(ctx context.Context, schoolID string, unitID *string, from, to *time.Time) ([]*entity.CalendarEvent, error)
	FindByID(ctx context.Context, id string) (*entity.CalendarEvent, error)
	Update(ctx context.Context, event *entity.CalendarEvent) error
	Delete(ctx context.Context, id string) error
}
