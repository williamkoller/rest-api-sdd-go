package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type InvoiceHandler struct {
	uc *usecase.InvoiceUseCase
}

func NewInvoiceHandler(uc *usecase.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{uc: uc}
}

func (h *InvoiceHandler) Generate(c *gin.Context) {
	var req struct {
		UnitID       string  `json:"unit_id" binding:"required"`
		AcademicYear int     `json:"academic_year" binding:"required"`
		Reference    string  `json:"reference" binding:"required"`
		DueDate      string  `json:"due_date" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid due_date")
		return
	}

	schoolID := middleware.GetSchoolID(c)
	count, err := h.uc.Generate(c.Request.Context(), req.UnitID, schoolID, req.AcademicYear, req.Reference, dueDate, req.Amount)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate invoices")
		return
	}
	response.JSON(c, http.StatusAccepted, gin.H{"queued": true, "estimatedCount": count})
}

func (h *InvoiceHandler) ListByStudent(c *gin.Context) {
	var year *int
	if v := c.Query("year"); v != "" {
		y, _ := strconv.Atoi(v) //nolint:errcheck
		year = &y
	}
	_ = year
	filters := repository.InvoiceFilters{Status: c.Query("status")}
	invoices, err := h.uc.GetByStudent(c.Request.Context(), c.Param("student_id"), filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch invoices")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapInvoices(invoices))
}

func (h *InvoiceHandler) Get(c *gin.Context) {
	invoice, err := h.uc.GetByID(c.Request.Context(), c.Param("invoice_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Invoice not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch invoice")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapInvoice(invoice))
}

func (h *InvoiceHandler) Pay(c *gin.Context) {
	var req struct {
		AmountPaid float64 `json:"amount_paid" binding:"required,gt=0"`
		Method     string  `json:"method" binding:"required"`
		GatewayRef string  `json:"gateway_ref"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	payment, err := h.uc.Pay(c.Request.Context(), c.Param("invoice_id"), repository.PaymentRequest{
		AmountPaid: req.AmountPaid,
		Method:     entity.PaymentMethod(req.Method),
		GatewayRef: req.GatewayRef,
		PaidAt:     time.Now(),
	})
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Invoice not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to record payment")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapPayment(payment))
}

func (h *InvoiceHandler) Receipt(c *gin.Context) {
	invoice, err := h.uc.GetByID(c.Request.Context(), c.Param("invoice_id"))
	if err != nil || invoice == nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Invoice not found")
		return
	}
	pdf := fmt.Sprintf("RECEIPT\nInvoice: %s\nStudent: %s\nAmount: R$ %.2f\nStatus: %s\nReference: %s\n",
		invoice.ID, invoice.StudentID, invoice.Amount, invoice.Status, invoice.Reference)
	c.Header("Content-Disposition", "attachment; filename=receipt-"+invoice.ID+".pdf")
	c.Data(http.StatusOK, "application/pdf", []byte(pdf))
}

func (h *InvoiceHandler) Delinquency(c *gin.Context) {
	daysOverdue, _ := strconv.Atoi(c.DefaultQuery("days_overdue", "5")) //nolint:errcheck
	var unitID *string
	if v := c.Query("unit_id"); v != "" {
		unitID = &v
	}
	entries, err := h.uc.GetDelinquency(c.Request.Context(), c.Param("school_id"), unitID, daysOverdue)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch delinquency report")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapDelinquencyEntries(entries))
}
