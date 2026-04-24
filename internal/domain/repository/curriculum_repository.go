package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type CurriculumRepository interface {
	BatchCreate(ctx context.Context, classID string, entries []*entity.CurriculumEntry) error
	FindByClassID(ctx context.Context, classID string) ([]*entity.CurriculumEntry, error)
	FindByID(ctx context.Context, id string) (*entity.CurriculumEntry, error)
	Update(ctx context.Context, entry *entity.CurriculumEntry) error
	Delete(ctx context.Context, id string) error
}
