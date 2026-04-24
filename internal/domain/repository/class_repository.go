package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ClassRepository interface {
	Create(ctx context.Context, class *entity.Class) error
	FindByID(ctx context.Context, id string) (*entity.Class, error)
	FindByUnitID(ctx context.Context, unitID string, year *int, active *bool) ([]*entity.Class, error)
	Update(ctx context.Context, class *entity.Class) error
	CountEnrollments(ctx context.Context, classID string) (int64, error)
}
