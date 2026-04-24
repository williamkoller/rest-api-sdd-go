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

type CalendarHandler struct {
	uc *usecase.CalendarUseCase
}

func NewCalendarHandler(uc *usecase.CalendarUseCase) *CalendarHandler {
	return &CalendarHandler{uc: uc}
}

func (h *CalendarHandler) Create(c *gin.Context) {
	var req struct {
		UnitID      string `json:"unit_id"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"required"`
		StartDate   string `json:"start_date" binding:"required"`
		EndDate     string `json:"end_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid start_date format")
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid end_date format")
		return
	}
	if end.Before(start) {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "end_date must be >= start_date")
		return
	}
	event := &entity.CalendarEvent{
		SchoolID:    c.Param("school_id"),
		UnitID:      req.UnitID,
		Title:       req.Title,
		Description: req.Description,
		Type:        entity.CalendarEventType(req.Type),
		StartDate:   start,
		EndDate:     end,
		CreatedBy:   middleware.GetUserID(c),
	}
	result, err := h.uc.Create(c.Request.Context(), event)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create calendar event")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapCalendarEvent(result))
}

func (h *CalendarHandler) List(c *gin.Context) {
	var from, to *time.Time
	var unitID *string
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
	if v := c.Query("unit_id"); v != "" {
		unitID = &v
	}
	events, err := h.uc.List(c.Request.Context(), c.Param("school_id"), unitID, from, to)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch calendar")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapCalendarEvents(events))
}

func (h *CalendarHandler) Update(c *gin.Context) {
	var req struct {
		UnitID      string `json:"unit_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Type        string `json:"type"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	start, _ := time.Parse("2006-01-02", req.StartDate) //nolint:errcheck
	end, _ := time.Parse("2006-01-02", req.EndDate)     //nolint:errcheck
	if !end.IsZero() && !start.IsZero() && end.Before(start) {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "end_date must be >= start_date")
		return
	}
	updates := &entity.CalendarEvent{
		UnitID:      req.UnitID,
		Title:       req.Title,
		Description: req.Description,
		Type:        entity.CalendarEventType(req.Type),
		StartDate:   start,
		EndDate:     end,
	}
	result, err := h.uc.Update(c.Request.Context(), c.Param("calendar_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Calendar event not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update calendar event")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapCalendarEvent(result))
}

func (h *CalendarHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("calendar_id")); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Calendar event not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete calendar event")
		return
	}
	response.JSON(c, http.StatusNoContent, nil)
}
