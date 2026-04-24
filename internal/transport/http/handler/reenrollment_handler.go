package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"

	"github.com/gin-gonic/gin"
)

type ReenrollmentHandler struct {
	uc *usecase.ReenrollmentUseCase
}

func NewReenrollmentHandler(uc *usecase.ReenrollmentUseCase) *ReenrollmentHandler {
	return &ReenrollmentHandler{uc: uc}
}

func (h *ReenrollmentHandler) CreateCampaign(c *gin.Context) {
	var req struct {
		UnitID       string `json:"unit_id"`
		AcademicYear int    `json:"academic_year" binding:"required"`
		Deadline     string `json:"deadline" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid deadline format, expected YYYY-MM-DD")
		return
	}
	schoolID := middleware.GetSchoolID(c)
	campaign, err := h.uc.CreateCampaign(c.Request.Context(), schoolID, req.UnitID, req.AcademicYear, deadline)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create campaign")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapReenrollmentCampaign(campaign))
}

func (h *ReenrollmentHandler) Dashboard(c *gin.Context) {
	dash, err := h.uc.GetDashboard(c.Request.Context(), c.Param("campaign_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Campaign not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch dashboard")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapCampaignDashboard(dash))
}

func (h *ReenrollmentHandler) Respond(c *gin.Context) {
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		Status    string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	record, err := h.uc.Respond(c.Request.Context(), c.Param("campaign_id"), req.StudentID, entity.ReenrollmentStatus(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Campaign not found")
		case errors.Is(err, usecase.ErrCampaignNotOpen):
			response.Error(c, http.StatusBadRequest, "CAMPAIGN_CLOSED", "Campaign is not open")
		case errors.Is(err, usecase.ErrDeadlinePassed):
			response.Error(c, http.StatusBadRequest, "DEADLINE_PASSED", "Reenrollment deadline has passed")
		case errors.Is(err, usecase.ErrOutstandingDebt):
			response.Error(c, http.StatusConflict, "OUTSTANDING_DEBT", "Student has outstanding invoices")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to record response")
		}
		return
	}
	response.JSON(c, http.StatusOK, dto.MapReenrollment(record))
}
