package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type InvoiceResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"studentId"`
	SchoolID  string    `json:"schoolId"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"dueDate"`
	Reference string    `json:"reference"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PaymentResponse struct {
	ID         string    `json:"id"`
	InvoiceID  string    `json:"invoiceId"`
	AmountPaid float64   `json:"amountPaid"`
	Method     string    `json:"method"`
	GatewayRef string    `json:"gatewayRef"`
	PaidAt     time.Time `json:"paidAt"`
}

func MapInvoice(e *entity.Invoice) *InvoiceResponse {
	return &InvoiceResponse{
		ID:        e.ID,
		StudentID: e.StudentID,
		SchoolID:  e.SchoolID,
		Amount:    e.Amount,
		DueDate:   e.DueDate,
		Reference: e.Reference,
		Status:    string(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func MapInvoices(list []*entity.Invoice) []*InvoiceResponse {
	res := make([]*InvoiceResponse, len(list))
	for i, e := range list {
		res[i] = MapInvoice(e)
	}
	return res
}

func MapPayment(e *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:         e.ID,
		InvoiceID:  e.InvoiceID,
		AmountPaid: e.AmountPaid,
		Method:     string(e.Method),
		GatewayRef: e.GatewayRef,
		PaidAt:     e.PaidAt,
	}
}

type DelinquencyEntryResponse struct {
	Student     *StudentResponse `json:"student"`
	Invoice     *InvoiceResponse `json:"invoice"`
	DaysOverdue int              `json:"daysOverdue"`
}

func MapDelinquencyEntry(e *repository.DelinquencyEntry) *DelinquencyEntryResponse {
	return &DelinquencyEntryResponse{
		Student:     MapStudent(e.Student),
		Invoice:     MapInvoice(e.Invoice),
		DaysOverdue: e.DaysOverdue,
	}
}

func MapDelinquencyEntries(list []*repository.DelinquencyEntry) []*DelinquencyEntryResponse {
	res := make([]*DelinquencyEntryResponse, len(list))
	for i, e := range list {
		res[i] = MapDelinquencyEntry(e)
	}
	return res
}
