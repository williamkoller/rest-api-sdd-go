package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type UnitRepository interface {
	Create(ctx context.Context, unit *entity.Unit) error
	FindByID(ctx context.Context, id string) (*entity.Unit, error)
	FindBySchoolID(ctx context.Context, schoolID string, active *bool) ([]*entity.Unit, error)
	Update(ctx context.Context, unit *entity.Unit) error
}
