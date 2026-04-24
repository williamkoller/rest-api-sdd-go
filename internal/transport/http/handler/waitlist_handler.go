package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type WaitlistHandler struct {
	uc *usecase.WaitlistUseCase
}

func NewWaitlistHandler(uc *usecase.WaitlistUseCase) *WaitlistHandler {
	return &WaitlistHandler{uc: uc}
}

func (h *WaitlistHandler) Register(c *gin.Context) {
	var req struct {
		GuardianName  string `json:"guardian_name" binding:"required"`
		GuardianEmail string `json:"guardian_email" binding:"required,email"`
		StudentName   string `json:"student_name" binding:"required"`
		GradeLevel    string `json:"grade_level" binding:"required"`
		AcademicYear  int    `json:"academic_year" binding:"required"`
		ReferralID    string `json:"referral_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	entry := &entity.WaitlistEntry{
		UnitID:        c.Param("unit_id"),
		GuardianName:  req.GuardianName,
		GuardianEmail: req.GuardianEmail,
		StudentName:   req.StudentName,
		GradeLevel:    req.GradeLevel,
		AcademicYear:  req.AcademicYear,
		ReferralID:    req.ReferralID,
	}
	result, err := h.uc.Register(c.Request.Context(), entry)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register on waitlist")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapWaitlistEntry(result))
}

func (h *WaitlistHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))         //nolint:errcheck
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20")) //nolint:errcheck
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	filters := repository.WaitlistFilters{
		Status:     c.Query("status"),
		GradeLevel: c.Query("grade_level"),
	}
	entries, total, err := h.uc.List(c.Request.Context(), c.Param("unit_id"), filters, page, perPage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch waitlist")
		return
	}
	response.Paginated(c, dto.MapWaitlistEntries(entries), page, perPage, total)
}

func (h *WaitlistHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	err := h.uc.UpdateStatus(c.Request.Context(), c.Param("waitlist_id"), entity.WaitlistStatus(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Waitlist entry not found")
		case errors.Is(err, usecase.ErrInvalidStatusTransition):
			response.Error(c, http.StatusBadRequest, "INVALID_TRANSITION", "Status transition not allowed")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update status")
		}
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": true})
}
