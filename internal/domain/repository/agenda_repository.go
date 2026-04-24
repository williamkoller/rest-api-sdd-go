package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AgendaRepository interface {
	Create(ctx context.Context, item *entity.AgendaItem) error
	FindByClassID(ctx context.Context, classID, typeFilter string, from, to *time.Time) ([]*entity.AgendaItem, error)
	FindByID(ctx context.Context, id string) (*entity.AgendaItem, error)
	Update(ctx context.Context, item *entity.AgendaItem) error
	Delete(ctx context.Context, id string) error
}
