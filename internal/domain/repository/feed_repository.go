package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type FeedRepository interface {
	Create(ctx context.Context, post *entity.FeedPost) error
	FindBySchoolID(ctx context.Context, schoolID string, unitID *string, page, perPage int) ([]*entity.FeedPost, int64, error)
	FindByID(ctx context.Context, id string) (*entity.FeedPost, error)
	Delete(ctx context.Context, id string) error
}
