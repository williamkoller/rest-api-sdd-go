package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domainrepo.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(ctx context.Context, menu *entity.Menu, items []*entity.MenuItem) error {
	menuID := uuid.New().String()
	ms := &model.MenuModel{
		ID:        menuID,
		UnitID:    menu.UnitID,
		WeekStart: menu.WeekStart,
		CreatedBy: menu.CreatedBy,
	}
	itemModels := make([]model.MenuItemModel, len(items))
	for i, item := range items {
		itemModels[i] = model.MenuItemModel{
			ID:          uuid.New().String(),
			MenuID:      menuID,
			DayOfWeek:   string(item.DayOfWeek),
			MealType:    string(item.MealType),
			Description: item.Description,
		}
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ms).Error; err != nil {
			return err
		}
		if len(itemModels) > 0 {
			return tx.Create(&itemModels).Error
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("menu repository: create: %w", err)
	}
	menu.ID = menuID
	return nil
}

func (r *menuRepository) FindByUnitAndWeek(ctx context.Context, unitID string, weekStart time.Time) (*entity.Menu, error) {
	var m model.MenuModel
	if err := r.db.WithContext(ctx).Preload("Items").Where("unit_id = ? AND week_start = ?", unitID, weekStart).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("menu repository: find by unit and week: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *menuRepository) FindByID(ctx context.Context, id string) (*entity.Menu, error) {
	var m model.MenuModel
	if err := r.db.WithContext(ctx).Preload("Items").First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("menu repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *menuRepository) ReplaceItems(ctx context.Context, menuID string, items []*entity.MenuItem) error {
	itemModels := make([]model.MenuItemModel, len(items))
	for i, item := range items {
		itemModels[i] = model.MenuItemModel{
			ID:          uuid.New().String(),
			MenuID:      menuID,
			DayOfWeek:   string(item.DayOfWeek),
			MealType:    string(item.MealType),
			Description: item.Description,
		}
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.MenuItemModel{}, "menu_id = ?", menuID).Error; err != nil {
			return err
		}
		if len(itemModels) > 0 {
			return tx.Create(&itemModels).Error
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("menu repository: replace items: %w", err)
	}
	return nil
}
