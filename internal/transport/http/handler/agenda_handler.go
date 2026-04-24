package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type AgendaHandler struct {
	uc *usecase.AgendaUseCase
}

func NewAgendaHandler(uc *usecase.AgendaUseCase) *AgendaHandler {
	return &AgendaHandler{uc: uc}
}

func (h *AgendaHandler) Create(c *gin.Context) {
	var req struct {
		Type        string  `json:"type" binding:"required"`
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid due_date format")
			return
		}
		dueDate = &t
	}
	item := &entity.AgendaItem{
		ClassID:     c.Param("class_id"),
		CreatedBy:   middleware.GetUserID(c),
		Type:        entity.AgendaItemType(req.Type),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}
	result, err := h.uc.Create(c.Request.Context(), item, middleware.GetRole(c))
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Not assigned to this class")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create agenda item")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapAgendaItem(result))
}

func (h *AgendaHandler) List(c *gin.Context) {
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
	items, err := h.uc.List(c.Request.Context(), c.Param("class_id"), c.Query("type"), from, to)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch agenda")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapAgendaItems(items))
}

func (h *AgendaHandler) Update(c *gin.Context) {
	var req struct {
		Type        string  `json:"type"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid due_date format")
			return
		}
		dueDate = &t
	}
	updates := &entity.AgendaItem{
		Type:        entity.AgendaItemType(req.Type),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}
	result, err := h.uc.Update(c.Request.Context(), c.Param("agenda_id"), middleware.GetUserID(c), middleware.GetRole(c), updates)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Agenda item not found")
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Cannot edit this item")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update agenda item")
		}
		return
	}
	response.JSON(c, http.StatusOK, dto.MapAgendaItem(result))
}

func (h *AgendaHandler) Delete(c *gin.Context) {
	err := h.uc.Delete(c.Request.Context(), c.Param("agenda_id"), middleware.GetUserID(c), middleware.GetRole(c))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Agenda item not found")
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Cannot delete this item")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete agenda item")
		}
		return
	}
	response.JSON(c, http.StatusNoContent, nil)
}
