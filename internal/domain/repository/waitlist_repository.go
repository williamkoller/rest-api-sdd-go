package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type WaitlistFilters struct {
	Status     string
	GradeLevel string
}

type WaitlistRepository interface {
	Create(ctx context.Context, entry *entity.WaitlistEntry) error
	FindByUnitID(ctx context.Context, unitID string, filters WaitlistFilters, page, perPage int) ([]*entity.WaitlistEntry, int64, error)
	FindByID(ctx context.Context, id string) (*entity.WaitlistEntry, error)
	UpdateStatus(ctx context.Context, id string, status entity.WaitlistStatus) error
}
