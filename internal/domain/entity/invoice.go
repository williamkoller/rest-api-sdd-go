package entity

import "time"

type InvoiceStatus string

const (
	InvoicePending   InvoiceStatus = "pending"
	InvoicePaid      InvoiceStatus = "paid"
	InvoiceOverdue   InvoiceStatus = "overdue"
	InvoiceCancelled InvoiceStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentPix        PaymentMethod = "pix"
	PaymentBoleto     PaymentMethod = "boleto"
	PaymentCreditCard PaymentMethod = "credit_card"
	PaymentManual     PaymentMethod = "manual"
)

type Invoice struct {
	ID        string
	StudentID string
	SchoolID  string
	Amount    float64
	DueDate   time.Time
	Reference string
	Status    InvoiceStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payment struct {
	ID         string
	InvoiceID  string
	AmountPaid float64
	PaidAt     time.Time
	Method     PaymentMethod
	GatewayRef string
	CreatedAt  time.Time
}
