package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var ErrInvalidStatusTransition = errors.New("invalid status transition")

var allowedTransitions = map[entity.WaitlistStatus][]entity.WaitlistStatus{
	entity.WaitlistWaiting:   {entity.WaitlistOfferMade, entity.WaitlistExpired},
	entity.WaitlistOfferMade: {entity.WaitlistAccepted, entity.WaitlistDeclined, entity.WaitlistExpired},
}

type WaitlistUseCase struct {
	repo repository.WaitlistRepository
}

func NewWaitlistUseCase(repo repository.WaitlistRepository) *WaitlistUseCase {
	return &WaitlistUseCase{repo: repo}
}

func (uc *WaitlistUseCase) Register(ctx context.Context, entry *entity.WaitlistEntry) (*entity.WaitlistEntry, error) {
	if err := uc.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("waitlist usecase: register: %w", err)
	}
	return entry, nil
}

func (uc *WaitlistUseCase) UpdateStatus(ctx context.Context, id string, newStatus entity.WaitlistStatus) error {
	entry, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("waitlist usecase: find for status update: %w", err)
	}
	if entry == nil {
		return ErrNotFound
	}
	allowed := allowedTransitions[entry.Status]
	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidStatusTransition
	}
	if err := uc.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return fmt.Errorf("waitlist usecase: update status: %w", err)
	}
	return nil
}

func (uc *WaitlistUseCase) List(ctx context.Context, unitID string, filters repository.WaitlistFilters, page, perPage int) ([]*entity.WaitlistEntry, int64, error) {
	entries, total, err := uc.repo.FindByUnitID(ctx, unitID, filters, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("waitlist usecase: list: %w", err)
	}
	return entries, total, nil
}
