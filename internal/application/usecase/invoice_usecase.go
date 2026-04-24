package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type InvoiceUseCase struct {
	invoiceRepo repository.InvoiceRepository
	studentRepo repository.StudentRepository
}

func NewInvoiceUseCase(invoiceRepo repository.InvoiceRepository, studentRepo repository.StudentRepository) *InvoiceUseCase {
	return &InvoiceUseCase{invoiceRepo: invoiceRepo, studentRepo: studentRepo}
}

func (uc *InvoiceUseCase) Generate(ctx context.Context, unitID, schoolID string, academicYear int, reference string, dueDate time.Time, amount float64) (int, error) {
	studentIDs, err := uc.invoiceRepo.CountActiveEnrollmentsByUnit(ctx, unitID)
	if err != nil {
		return 0, fmt.Errorf("invoice usecase: generate: %w", err)
	}

	invoices := make([]*entity.Invoice, 0, len(studentIDs))
	for _, sid := range studentIDs {
		invoices = append(invoices, &entity.Invoice{
			ID:        uuid.New().String(),
			StudentID: sid,
			SchoolID:  schoolID,
			Amount:    amount,
			DueDate:   dueDate,
			Reference: reference,
			Status:    entity.InvoicePending,
		})
	}

	if len(invoices) == 0 {
		return 0, nil
	}

	if err := uc.invoiceRepo.BatchCreate(ctx, invoices); err != nil {
		return 0, fmt.Errorf("invoice usecase: batch create: %w", err)
	}
	return len(invoices), nil
}

func (uc *InvoiceUseCase) GetByStudent(ctx context.Context, studentID string, filters repository.InvoiceFilters) ([]*entity.Invoice, error) {
	invoices, err := uc.invoiceRepo.FindByStudentID(ctx, studentID, filters)
	if err != nil {
		return nil, fmt.Errorf("invoice usecase: get by student: %w", err)
	}
	return invoices, nil
}

func (uc *InvoiceUseCase) GetByID(ctx context.Context, id string) (*entity.Invoice, error) {
	invoice, err := uc.invoiceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invoice usecase: get by id: %w", err)
	}
	if invoice == nil {
		return nil, ErrNotFound
	}
	return invoice, nil
}

func (uc *InvoiceUseCase) Pay(ctx context.Context, invoiceID string, req repository.PaymentRequest) (*entity.Payment, error) {
	invoice, err := uc.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice usecase: pay find: %w", err)
	}
	if invoice == nil {
		return nil, ErrNotFound
	}

	payment := &entity.Payment{
		InvoiceID:  invoiceID,
		AmountPaid: req.AmountPaid,
		PaidAt:     req.PaidAt,
		Method:     req.Method,
		GatewayRef: req.GatewayRef,
	}
	if err := uc.invoiceRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("invoice usecase: create payment: %w", err)
	}

	if err := uc.invoiceRepo.UpdateStatus(ctx, invoiceID, entity.InvoicePaid); err != nil {
		return nil, fmt.Errorf("invoice usecase: update status: %w", err)
	}

	return payment, nil
}

func (uc *InvoiceUseCase) GetDelinquency(ctx context.Context, schoolID string, unitID *string, daysOverdue int) ([]*repository.DelinquencyEntry, error) {
	entries, err := uc.invoiceRepo.FindOverdue(ctx, schoolID, unitID, daysOverdue)
	if err != nil {
		return nil, fmt.Errorf("invoice usecase: delinquency: %w", err)
	}
	return entries, nil
}
