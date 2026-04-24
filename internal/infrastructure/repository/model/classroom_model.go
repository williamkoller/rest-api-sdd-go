package model

import "github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"

type ClassroomModel struct {
	ID       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UnitID   string `gorm:"type:uuid;not null;index"`
	Code     string `gorm:"not null;size:20"`
	Capacity int    `gorm:"not null"`
	Active   bool   `gorm:"not null;default:true"`
}

func (ClassroomModel) TableName() string { return "classrooms" }

func (m *ClassroomModel) ToEntity() *entity.Classroom {
	return &entity.Classroom{
		ID:       m.ID,
		UnitID:   m.UnitID,
		Code:     m.Code,
		Capacity: m.Capacity,
		Active:   m.Active,
	}
}

func ClassroomFromEntity(c *entity.Classroom) *ClassroomModel {
	return &ClassroomModel{
		ID:       c.ID,
		UnitID:   c.UnitID,
		Code:     c.Code,
		Capacity: c.Capacity,
		Active:   c.Active,
	}
}
