package usecase

import (
	"context"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type CurriculumUseCase struct {
	repo repository.CurriculumRepository
}

func NewCurriculumUseCase(repo repository.CurriculumRepository) *CurriculumUseCase {
	return &CurriculumUseCase{repo: repo}
}

func (uc *CurriculumUseCase) SetCurriculum(ctx context.Context, classID string, entries []*entity.CurriculumEntry, role string) error {
	if role != string(entity.RoleSchoolAdmin) && role != string(entity.RoleSuperAdmin) {
		return ErrForbidden
	}
	if err := uc.repo.BatchCreate(ctx, classID, entries); err != nil {
		return fmt.Errorf("curriculum usecase: set curriculum: %w", err)
	}
	return nil
}

func (uc *CurriculumUseCase) GetCurriculum(ctx context.Context, classID string) ([]*entity.CurriculumEntry, error) {
	entries, err := uc.repo.FindByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("curriculum usecase: get curriculum: %w", err)
	}
	return entries, nil
}

func (uc *CurriculumUseCase) UpdateEntry(ctx context.Context, id string, entry *entity.CurriculumEntry) error {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("curriculum usecase: update find: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}
	entry.ID = id
	if err := uc.repo.Update(ctx, entry); err != nil {
		return fmt.Errorf("curriculum usecase: update: %w", err)
	}
	return nil
}

func (uc *CurriculumUseCase) DeleteEntry(ctx context.Context, id string) error {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("curriculum usecase: delete find: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("curriculum usecase: delete: %w", err)
	}
	return nil
}
