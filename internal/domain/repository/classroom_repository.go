package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ClassroomRepository interface {
	Create(ctx context.Context, classroom *entity.Classroom) error
	FindByID(ctx context.Context, id string) (*entity.Classroom, error)
	FindByUnitID(ctx context.Context, unitID string) ([]*entity.Classroom, error)
	Update(ctx context.Context, classroom *entity.Classroom) error
}
