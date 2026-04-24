package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type MenuModel struct {
	ID        string          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UnitID    string          `gorm:"type:uuid;not null;index"`
	WeekStart time.Time       `gorm:"type:date;not null"`
	CreatedBy string          `gorm:"type:uuid;not null"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	Items     []MenuItemModel `gorm:"foreignKey:MenuID"`
}

func (MenuModel) TableName() string { return "menus" }

func (m *MenuModel) ToEntity() *entity.Menu {
	items := make([]*entity.MenuItem, len(m.Items))
	for i, item := range m.Items {
		item := item
		items[i] = item.ToEntity()
	}
	return &entity.Menu{
		ID:        m.ID,
		UnitID:    m.UnitID,
		WeekStart: m.WeekStart,
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Items:     items,
	}
}

type MenuItemModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MenuID      string `gorm:"type:uuid;not null;index"`
	DayOfWeek   string `gorm:"not null;size:10"`
	MealType    string `gorm:"not null;size:15"`
	Description string `gorm:"type:text;not null"`
}

func (MenuItemModel) TableName() string { return "menu_items" }

func (m *MenuItemModel) ToEntity() *entity.MenuItem {
	return &entity.MenuItem{
		ID:          m.ID,
		MenuID:      m.MenuID,
		DayOfWeek:   entity.DayOfWeek(m.DayOfWeek),
		MealType:    entity.MealType(m.MealType),
		Description: m.Description,
	}
}
