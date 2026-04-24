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

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) domainrepo.InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) BatchCreate(ctx context.Context, invoices []*entity.Invoice) error {
	ms := make([]model.InvoiceModel, len(invoices))
	for i, inv := range invoices {
		inv.ID = uuid.New().String()
		ms[i] = *model.InvoiceFromEntity(inv)
	}
	if err := r.db.WithContext(ctx).Create(&ms).Error; err != nil {
		return fmt.Errorf("invoice repository: batch create: %w", err)
	}
	return nil
}

func (r *invoiceRepository) FindByStudentID(ctx context.Context, studentID string, filters domainrepo.InvoiceFilters) ([]*entity.Invoice, error) {
	var ms []model.InvoiceModel
	q := r.db.WithContext(ctx).Where("student_id = ?", studentID)
	if filters.Status != "" {
		q = q.Where("status = ?", filters.Status)
	}
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("invoice repository: find by student: %w", err)
	}
	result := make([]*entity.Invoice, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *invoiceRepository) FindByID(ctx context.Context, id string) (*entity.Invoice, error) {
	var m model.InvoiceModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("invoice repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id string, status entity.InvoiceStatus) error {
	err := r.db.WithContext(ctx).Model(&model.InvoiceModel{}).
		Where("id = ?", id).
		Updates(map[string]any{"status": string(status), "updated_at": time.Now()}).Error
	if err != nil {
		return fmt.Errorf("invoice repository: update status: %w", err)
	}
	return nil
}

func (r *invoiceRepository) HasOverdueInvoices(ctx context.Context, studentID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.InvoiceModel{}).
		Where("student_id = ? AND status = 'overdue'", studentID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("invoice repository: has overdue: %w", err)
	}
	return count > 0, nil
}

func (r *invoiceRepository) FindOverdue(ctx context.Context, schoolID string, unitID *string, daysOverdue int) ([]*domainrepo.DelinquencyEntry, error) {
	cutoff := time.Now().AddDate(0, 0, -daysOverdue)
	type row struct {
		model.InvoiceModel
		StudentName string
	}
	var rows []row
	q := r.db.WithContext(ctx).
		Table("invoices").
		Select("invoices.*, students.name as student_name").
		Joins("JOIN students ON students.id = invoices.student_id").
		Where("invoices.school_id = ? AND invoices.status IN ('pending','overdue') AND invoices.due_date <= ?", schoolID, cutoff)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("invoice repository: find overdue: %w", err)
	}
	result := make([]*domainrepo.DelinquencyEntry, len(rows))
	for i, row := range rows {
		days := int(time.Since(row.DueDate).Hours() / 24)
		result[i] = &domainrepo.DelinquencyEntry{
			Student:     &entity.Student{ID: row.StudentID, Name: row.StudentName},
			Invoice:     row.ToEntity(),
			DaysOverdue: days,
		}
	}
	return result, nil
}

func (r *invoiceRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	m := &model.PaymentModel{
		ID:         uuid.New().String(),
		InvoiceID:  payment.InvoiceID,
		AmountPaid: payment.AmountPaid,
		PaidAt:     payment.PaidAt,
		Method:     string(payment.Method),
		GatewayRef: payment.GatewayRef,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("invoice repository: create payment: %w", err)
	}
	payment.ID = m.ID
	return nil
}

func (r *invoiceRepository) CountActiveEnrollmentsByUnit(ctx context.Context, unitID string) ([]string, error) {
	var studentIDs []string
	err := r.db.WithContext(ctx).
		Table("enrollments").
		Select("DISTINCT enrollments.student_id").
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.unit_id = ? AND enrollments.status = 'active'", unitID).
		Pluck("enrollments.student_id", &studentIDs).Error
	if err != nil {
		return nil, fmt.Errorf("invoice repository: count enrollments by unit: %w", err)
	}
	return studentIDs, nil
}
