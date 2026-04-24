package handler

import (
	"errors"
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

type TicketHandler struct {
	uc *usecase.TicketUseCase
}

func NewTicketHandler(uc *usecase.TicketUseCase) *TicketHandler {
	return &TicketHandler{uc: uc}
}

func (h *TicketHandler) Create(c *gin.Context) {
	var req struct {
		UnitID   string `json:"unit_id"`
		Category string `json:"category" binding:"required"`
		Subject  string `json:"subject" binding:"required"`
		Message  string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	ticket := &entity.Ticket{
		SchoolID:    middleware.GetSchoolID(c),
		UnitID:      req.UnitID,
		RequesterID: middleware.GetUserID(c),
		Category:    entity.TicketCategory(req.Category),
		Subject:     req.Subject,
	}
	result, err := h.uc.Create(c.Request.Context(), ticket, req.Message)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create ticket")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapTicket(result))
}

func (h *TicketHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))         //nolint:errcheck
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20")) //nolint:errcheck
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	filters := repository.TicketFilters{
		Status:   c.Query("status"),
		Category: c.Query("category"),
	}
	tickets, total, err := h.uc.List(c.Request.Context(), middleware.GetSchoolID(c), middleware.GetUserID(c), middleware.GetRole(c), filters, page, perPage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tickets")
		return
	}
	response.Paginated(c, dto.MapTickets(tickets), page, perPage, total)
}

func (h *TicketHandler) Get(c *gin.Context) {
	ticket, err := h.uc.GetByID(c.Request.Context(), c.Param("ticket_id"), middleware.GetUserID(c), middleware.GetRole(c))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Ticket not found")
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Cannot access this ticket")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch ticket")
		}
		return
	}
	response.JSON(c, http.StatusOK, dto.MapTicket(ticket))
}

func (h *TicketHandler) Reply(c *gin.Context) {
	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	msg, err := h.uc.Reply(c.Request.Context(), c.Param("ticket_id"), req.Body, middleware.GetUserID(c), middleware.GetRole(c))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Ticket not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to send reply")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapTicketMessage(msg))
}

func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	err := h.uc.UpdateStatus(c.Request.Context(), c.Param("ticket_id"), entity.TicketStatus(req.Status), middleware.GetRole(c))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Ticket not found")
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Only staff can update ticket status")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update status")
		}
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": true})
}

func (h *TicketHandler) Report(c *gin.Context) {
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = &t
		}
	}
	report, err := h.uc.GetReport(c.Request.Context(), c.Param("school_id"), from, to)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate report")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapTicketReport(report))
}
