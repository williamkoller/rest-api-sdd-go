package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type ClassroomUseCase struct {
	repo repository.ClassroomRepository
}

func NewClassroomUseCase(repo repository.ClassroomRepository) *ClassroomUseCase {
	return &ClassroomUseCase{repo: repo}
}

func (uc *ClassroomUseCase) Create(ctx context.Context, classroom *entity.Classroom) error {
	classroom.ID = uuid.New().String()
	classroom.Active = true
	if err := uc.repo.Create(ctx, classroom); err != nil {
		return fmt.Errorf("classroom usecase: create: %w", err)
	}
	return nil
}

func (uc *ClassroomUseCase) FindByUnitID(ctx context.Context, unitID string) ([]*entity.Classroom, error) {
	classrooms, err := uc.repo.FindByUnitID(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("classroom usecase: find by unit id: %w", err)
	}
	return classrooms, nil
}

func (uc *ClassroomUseCase) FindByID(ctx context.Context, id string) (*entity.Classroom, error) {
	classroom, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("classroom usecase: find by id: %w", err)
	}
	if classroom == nil {
		return nil, ErrNotFound
	}
	return classroom, nil
}

func (uc *ClassroomUseCase) Update(ctx context.Context, id string, updates map[string]any) (*entity.Classroom, error) {
	classroom, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("classroom usecase: update: %w", err)
	}
	if classroom == nil {
		return nil, ErrNotFound
	}
	if v, ok := updates["code"].(string); ok {
		classroom.Code = v
	}
	if v, ok := updates["capacity"].(float64); ok {
		classroom.Capacity = int(v)
	}
	if v, ok := updates["active"].(bool); ok {
		classroom.Active = v
	}
	if err := uc.repo.Update(ctx, classroom); err != nil {
		return nil, fmt.Errorf("classroom usecase: update: %w", err)
	}
	return classroom, nil
}
