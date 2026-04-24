package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type ClassWithCount struct {
	Class           *entity.Class
	EnrollmentCount int64
}

type ClassUseCase struct {
	repo repository.ClassRepository
}

func NewClassUseCase(repo repository.ClassRepository) *ClassUseCase {
	return &ClassUseCase{repo: repo}
}

func (uc *ClassUseCase) Create(ctx context.Context, class *entity.Class) error {
	class.ID = uuid.New().String()
	class.Active = true
	if err := uc.repo.Create(ctx, class); err != nil {
		return fmt.Errorf("class usecase: create: %w", err)
	}
	return nil
}

func (uc *ClassUseCase) FindByUnitID(ctx context.Context, unitID string, year *int, active *bool) ([]*ClassWithCount, error) {
	classes, err := uc.repo.FindByUnitID(ctx, unitID, year, active)
	if err != nil {
		return nil, fmt.Errorf("class usecase: find by unit id: %w", err)
	}
	result := make([]*ClassWithCount, len(classes))
	for i, c := range classes {
		count, _ := uc.repo.CountEnrollments(ctx, c.ID) //nolint:errcheck
		result[i] = &ClassWithCount{Class: c, EnrollmentCount: count}
	}
	return result, nil
}

func (uc *ClassUseCase) FindByID(ctx context.Context, id string) (*ClassWithCount, error) {
	class, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("class usecase: find by id: %w", err)
	}
	if class == nil {
		return nil, ErrNotFound
	}
	count, _ := uc.repo.CountEnrollments(ctx, id) //nolint:errcheck
	return &ClassWithCount{Class: class, EnrollmentCount: count}, nil
}

func (uc *ClassUseCase) Update(ctx context.Context, id string, updates map[string]any) (*entity.Class, error) {
	class, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("class usecase: update: %w", err)
	}
	if class == nil {
		return nil, ErrNotFound
	}
	if v, ok := updates["name"].(string); ok {
		class.Name = v
	}
	if v, ok := updates["classroom_id"].(string); ok {
		class.ClassroomID = v
	}
	if v, ok := updates["active"].(bool); ok {
		class.Active = v
	}
	if err := uc.repo.Update(ctx, class); err != nil {
		return nil, fmt.Errorf("class usecase: update: %w", err)
	}
	return class, nil
}
