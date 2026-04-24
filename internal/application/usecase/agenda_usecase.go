package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type AgendaUseCase struct {
	agendaRepo repository.AgendaRepository
	userRepo   repository.UserRepository
}

func NewAgendaUseCase(agendaRepo repository.AgendaRepository, userRepo repository.UserRepository) *AgendaUseCase {
	return &AgendaUseCase{agendaRepo: agendaRepo, userRepo: userRepo}
}

func (uc *AgendaUseCase) Create(ctx context.Context, item *entity.AgendaItem, role string) (*entity.AgendaItem, error) {
	if role == string(entity.RoleTeacher) {
		isTeacher, err := uc.userRepo.IsTeacherOfClass(ctx, item.CreatedBy, item.ClassID)
		if err != nil {
			return nil, fmt.Errorf("agenda usecase: check teacher: %w", err)
		}
		if !isTeacher {
			return nil, ErrForbidden
		}
	}
	if err := uc.agendaRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("agenda usecase: create: %w", err)
	}
	return item, nil
}

func (uc *AgendaUseCase) List(ctx context.Context, classID, typeFilter string, from, to *time.Time) ([]*entity.AgendaItem, error) {
	items, err := uc.agendaRepo.FindByClassID(ctx, classID, typeFilter, from, to)
	if err != nil {
		return nil, fmt.Errorf("agenda usecase: list: %w", err)
	}
	return items, nil
}

func (uc *AgendaUseCase) Update(ctx context.Context, id, userID, role string, updates *entity.AgendaItem) (*entity.AgendaItem, error) {
	existing, err := uc.agendaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("agenda usecase: update find: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	if role == string(entity.RoleTeacher) && existing.CreatedBy != userID {
		return nil, ErrForbidden
	}
	updates.ID = id
	if err := uc.agendaRepo.Update(ctx, updates); err != nil {
		return nil, fmt.Errorf("agenda usecase: update: %w", err)
	}
	return updates, nil
}

func (uc *AgendaUseCase) Delete(ctx context.Context, id, userID, role string) error {
	existing, err := uc.agendaRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("agenda usecase: delete find: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}
	if role == string(entity.RoleTeacher) && existing.CreatedBy != userID {
		return ErrForbidden
	}
	if err := uc.agendaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("agenda usecase: delete: %w", err)
	}
	return nil
}
