package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type SchoolRepository interface {
	Create(ctx context.Context, school *entity.School) error
	FindByID(ctx context.Context, id string) (*entity.School, error)
	FindAll(ctx context.Context, page, perPage int, active *bool) ([]*entity.School, int64, error)
	Update(ctx context.Context, school *entity.School) error
}
