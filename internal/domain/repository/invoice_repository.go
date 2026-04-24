package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type InvoiceFilters struct {
	Status string
	Year   *int
}

type DelinquencyEntry struct {
	Student     *entity.Student
	Invoice     *entity.Invoice
	DaysOverdue int
}

type InvoiceRepository interface {
	BatchCreate(ctx context.Context, invoices []*entity.Invoice) error
	FindByStudentID(ctx context.Context, studentID string, filters InvoiceFilters) ([]*entity.Invoice, error)
	FindByID(ctx context.Context, id string) (*entity.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status entity.InvoiceStatus) error
	HasOverdueInvoices(ctx context.Context, studentID string) (bool, error)
	FindOverdue(ctx context.Context, schoolID string, unitID *string, daysOverdue int) ([]*DelinquencyEntry, error)
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	CountActiveEnrollmentsByUnit(ctx context.Context, unitID string) ([]string, error)
}

type PaymentRequest struct {
	AmountPaid float64
	Method     entity.PaymentMethod
	GatewayRef string
	PaidAt     time.Time
}
