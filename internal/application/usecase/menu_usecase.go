package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var ErrMenuAlreadyExists = errors.New("menu already exists for this week")

type MenuUseCase struct {
	repo repository.MenuRepository
}

func NewMenuUseCase(repo repository.MenuRepository) *MenuUseCase {
	return &MenuUseCase{repo: repo}
}

func (uc *MenuUseCase) Publish(ctx context.Context, unitID string, weekStart time.Time, items []*entity.MenuItem, userID, role string) (*entity.Menu, error) {
	if !publisherRoles[role] {
		return nil, ErrForbidden
	}
	existing, err := uc.repo.FindByUnitAndWeek(ctx, unitID, weekStart)
	if err != nil {
		return nil, fmt.Errorf("menu usecase: check existing: %w", err)
	}
	if existing != nil {
		return nil, ErrMenuAlreadyExists
	}
	menu := &entity.Menu{
		UnitID:    unitID,
		WeekStart: weekStart,
		CreatedBy: userID,
	}
	if err := uc.repo.Create(ctx, menu, items); err != nil {
		return nil, fmt.Errorf("menu usecase: publish: %w", err)
	}
	menu.Items = items
	return menu, nil
}

func (uc *MenuUseCase) GetMenu(ctx context.Context, unitID string, weekStart time.Time) (*entity.Menu, error) {
	menu, err := uc.repo.FindByUnitAndWeek(ctx, unitID, weekStart)
	if err != nil {
		return nil, fmt.Errorf("menu usecase: get menu: %w", err)
	}
	if menu == nil {
		return nil, ErrNotFound
	}
	return menu, nil
}

func (uc *MenuUseCase) ReplaceMenu(ctx context.Context, menuID string, items []*entity.MenuItem) error {
	menu, err := uc.repo.FindByID(ctx, menuID)
	if err != nil {
		return fmt.Errorf("menu usecase: replace find: %w", err)
	}
	if menu == nil {
		return ErrNotFound
	}
	if err := uc.repo.ReplaceItems(ctx, menuID, items); err != nil {
		return fmt.Errorf("menu usecase: replace items: %w", err)
	}
	return nil
}
