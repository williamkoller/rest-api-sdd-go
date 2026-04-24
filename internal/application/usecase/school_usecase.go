package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var ErrNotFound = errors.New("resource not found")

type SchoolUseCase struct {
	repo repository.SchoolRepository
}

func NewSchoolUseCase(repo repository.SchoolRepository) *SchoolUseCase {
	return &SchoolUseCase{repo: repo}
}

func (uc *SchoolUseCase) Create(ctx context.Context, school *entity.School) error {
	school.ID = uuid.New().String()
	school.Active = true
	if err := uc.repo.Create(ctx, school); err != nil {
		return fmt.Errorf("school usecase: create: %w", err)
	}
	return nil
}

func (uc *SchoolUseCase) FindAll(ctx context.Context, page, perPage int, active *bool) ([]*entity.School, int64, error) {
	schools, total, err := uc.repo.FindAll(ctx, page, perPage, active)
	if err != nil {
		return nil, 0, fmt.Errorf("school usecase: find all: %w", err)
	}
	return schools, total, nil
}

func (uc *SchoolUseCase) FindByID(ctx context.Context, id string) (*entity.School, error) {
	school, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("school usecase: find by id: %w", err)
	}
	if school == nil {
		return nil, ErrNotFound
	}
	return school, nil
}

func (uc *SchoolUseCase) Update(ctx context.Context, id string, updates map[string]any) (*entity.School, error) {
	school, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("school usecase: update: %w", err)
	}
	if school == nil {
		return nil, ErrNotFound
	}
	if v, ok := updates["name"].(string); ok {
		school.Name = v
	}
	if v, ok := updates["email"].(string); ok {
		school.Email = v
	}
	if v, ok := updates["phone"].(string); ok {
		school.Phone = v
	}
	if v, ok := updates["active"].(bool); ok {
		school.Active = v
	}
	if err := uc.repo.Update(ctx, school); err != nil {
		return nil, fmt.Errorf("school usecase: update: %w", err)
	}
	return school, nil
}
