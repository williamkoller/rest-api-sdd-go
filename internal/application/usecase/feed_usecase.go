package usecase

import (
	"context"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var publisherRoles = map[string]bool{
	string(entity.RoleSchoolAdmin): true,
	string(entity.RoleUnitStaff):   true,
	string(entity.RoleSuperAdmin):  true,
}

type FeedUseCase struct {
	repo repository.FeedRepository
}

func NewFeedUseCase(repo repository.FeedRepository) *FeedUseCase {
	return &FeedUseCase{repo: repo}
}

func (uc *FeedUseCase) Publish(ctx context.Context, post *entity.FeedPost, role string) (*entity.FeedPost, error) {
	if !publisherRoles[role] {
		return nil, ErrForbidden
	}
	if err := uc.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("feed usecase: publish: %w", err)
	}
	return post, nil
}

func (uc *FeedUseCase) List(ctx context.Context, schoolID string, unitID *string, page, perPage int) ([]*entity.FeedPost, int64, error) {
	posts, total, err := uc.repo.FindBySchoolID(ctx, schoolID, unitID, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("feed usecase: list: %w", err)
	}
	return posts, total, nil
}

func (uc *FeedUseCase) Delete(ctx context.Context, id, userID, role string) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("feed usecase: delete find: %w", err)
	}
	if post == nil {
		return ErrNotFound
	}
	if !publisherRoles[role] && post.AuthorID != userID {
		return ErrForbidden
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("feed usecase: delete: %w", err)
	}
	return nil
}
