package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type CalendarUseCase struct {
	repo repository.CalendarRepository
}

func NewCalendarUseCase(repo repository.CalendarRepository) *CalendarUseCase {
	return &CalendarUseCase{repo: repo}
}

func (uc *CalendarUseCase) Create(ctx context.Context, event *entity.CalendarEvent) (*entity.CalendarEvent, error) {
	if err := uc.repo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("calendar usecase: create: %w", err)
	}
	return event, nil
}

func (uc *CalendarUseCase) List(ctx context.Context, schoolID string, unitID *string, from, to *time.Time) ([]*entity.CalendarEvent, error) {
	events, err := uc.repo.FindBySchoolID(ctx, schoolID, unitID, from, to)
	if err != nil {
		return nil, fmt.Errorf("calendar usecase: list: %w", err)
	}
	return events, nil
}

func (uc *CalendarUseCase) Update(ctx context.Context, id string, updates *entity.CalendarEvent) (*entity.CalendarEvent, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("calendar usecase: update find: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	updates.ID = id
	if err := uc.repo.Update(ctx, updates); err != nil {
		return nil, fmt.Errorf("calendar usecase: update: %w", err)
	}
	return updates, nil
}

func (uc *CalendarUseCase) Delete(ctx context.Context, id string) error {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("calendar usecase: delete find: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("calendar usecase: delete: %w", err)
	}
	return nil
}
