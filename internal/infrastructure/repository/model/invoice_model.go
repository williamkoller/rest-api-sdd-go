package model

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type InvoiceModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	StudentID string    `gorm:"type:uuid;not null;index"`
	SchoolID  string    `gorm:"type:uuid;not null;index"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	DueDate   time.Time `gorm:"not null"`
	Reference string    `gorm:"not null;size:20"`
	Status    string    `gorm:"not null;default:'pending'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (InvoiceModel) TableName() string { return "invoices" }

func (m *InvoiceModel) ToEntity() *entity.Invoice {
	return &entity.Invoice{
		ID:        m.ID,
		StudentID: m.StudentID,
		SchoolID:  m.SchoolID,
		Amount:    m.Amount,
		DueDate:   m.DueDate,
		Reference: m.Reference,
		Status:    entity.InvoiceStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func InvoiceFromEntity(i *entity.Invoice) *InvoiceModel {
	return &InvoiceModel{
		ID:        i.ID,
		StudentID: i.StudentID,
		SchoolID:  i.SchoolID,
		Amount:    i.Amount,
		DueDate:   i.DueDate,
		Reference: i.Reference,
		Status:    string(i.Status),
	}
}

type PaymentModel struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	InvoiceID  string    `gorm:"type:uuid;not null;index"`
	AmountPaid float64   `gorm:"type:decimal(10,2);not null"`
	PaidAt     time.Time `gorm:"not null"`
	Method     string    `gorm:"not null"`
	GatewayRef string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (PaymentModel) TableName() string { return "payments" }

func (m *PaymentModel) ToEntity() *entity.Payment {
	return &entity.Payment{
		ID:         m.ID,
		InvoiceID:  m.InvoiceID,
		AmountPaid: m.AmountPaid,
		PaidAt:     m.PaidAt,
		Method:     entity.PaymentMethod(m.Method),
		GatewayRef: m.GatewayRef,
		CreatedAt:  m.CreatedAt,
	}
}
