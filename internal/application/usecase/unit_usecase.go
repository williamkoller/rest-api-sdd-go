package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type UnitUseCase struct {
	repo repository.UnitRepository
}

func NewUnitUseCase(repo repository.UnitRepository) *UnitUseCase {
	return &UnitUseCase{repo: repo}
}

func (uc *UnitUseCase) Create(ctx context.Context, unit *entity.Unit) error {
	unit.ID = uuid.New().String()
	unit.Active = true
	if err := uc.repo.Create(ctx, unit); err != nil {
		return fmt.Errorf("unit usecase: create: %w", err)
	}
	return nil
}

func (uc *UnitUseCase) FindBySchoolID(ctx context.Context, schoolID string, active *bool) ([]*entity.Unit, error) {
	units, err := uc.repo.FindBySchoolID(ctx, schoolID, active)
	if err != nil {
		return nil, fmt.Errorf("unit usecase: find by school id: %w", err)
	}
	return units, nil
}

func (uc *UnitUseCase) FindByID(ctx context.Context, id string) (*entity.Unit, error) {
	unit, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unit usecase: find by id: %w", err)
	}
	if unit == nil {
		return nil, ErrNotFound
	}
	return unit, nil
}

func (uc *UnitUseCase) Update(ctx context.Context, id string, updates map[string]any) (*entity.Unit, error) {
	unit, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unit usecase: update: %w", err)
	}
	if unit == nil {
		return nil, ErrNotFound
	}
	if v, ok := updates["name"].(string); ok {
		unit.Name = v
	}
	if v, ok := updates["address"].(string); ok {
		unit.Address = v
	}
	if v, ok := updates["city"].(string); ok {
		unit.City = v
	}
	if v, ok := updates["state"].(string); ok {
		unit.State = v
	}
	if v, ok := updates["zip_code"].(string); ok {
		unit.ZipCode = v
	}
	if v, ok := updates["phone"].(string); ok {
		unit.Phone = v
	}
	if v, ok := updates["active"].(bool); ok {
		unit.Active = v
	}
	if err := uc.repo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("unit usecase: update: %w", err)
	}
	return unit, nil
}
