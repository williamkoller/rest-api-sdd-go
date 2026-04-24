package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type MenuRepository interface {
	Create(ctx context.Context, menu *entity.Menu, items []*entity.MenuItem) error
	FindByUnitAndWeek(ctx context.Context, unitID string, weekStart time.Time) (*entity.Menu, error)
	FindByID(ctx context.Context, id string) (*entity.Menu, error)
	ReplaceItems(ctx context.Context, menuID string, items []*entity.MenuItem) error
}
