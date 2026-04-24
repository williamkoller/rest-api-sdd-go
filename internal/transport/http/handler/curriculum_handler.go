package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type CurriculumHandler struct {
	uc *usecase.CurriculumUseCase
}

func NewCurriculumHandler(uc *usecase.CurriculumUseCase) *CurriculumHandler {
	return &CurriculumHandler{uc: uc}
}

type curriculumEntryRequest struct {
	Subject   string `json:"subject" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
	DayOfWeek string `json:"day_of_week" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

func (h *CurriculumHandler) BatchCreate(c *gin.Context) {
	var req struct {
		Entries []curriculumEntryRequest `json:"entries" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	entries := make([]*entity.CurriculumEntry, len(req.Entries))
	for i, e := range req.Entries {
		entries[i] = &entity.CurriculumEntry{
			Subject:   e.Subject,
			TeacherID: e.TeacherID,
			DayOfWeek: entity.DayOfWeek(e.DayOfWeek),
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
		}
	}
	if err := h.uc.SetCurriculum(c.Request.Context(), c.Param("class_id"), entries, middleware.GetRole(c)); err != nil {
		switch {
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Only school admin can set curriculum")
		default:
			response.Error(c, http.StatusConflict, "SCHEDULE_CONFLICT", "Time slot conflict detected")
		}
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": true})
}

func (h *CurriculumHandler) List(c *gin.Context) {
	entries, err := h.uc.GetCurriculum(c.Request.Context(), c.Param("class_id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch curriculum")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapCurriculumEntries(entries))
}

func (h *CurriculumHandler) Update(c *gin.Context) {
	var req curriculumEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	entry := &entity.CurriculumEntry{
		Subject:   req.Subject,
		TeacherID: req.TeacherID,
		DayOfWeek: entity.DayOfWeek(req.DayOfWeek),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if err := h.uc.UpdateEntry(c.Request.Context(), c.Param("curriculum_id"), entry); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Curriculum entry not found")
			return
		}
		response.Error(c, http.StatusConflict, "SCHEDULE_CONFLICT", "Time slot conflict detected")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": true})
}

func (h *CurriculumHandler) Delete(c *gin.Context) {
	if err := h.uc.DeleteEntry(c.Request.Context(), c.Param("curriculum_id")); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Curriculum entry not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete curriculum entry")
		return
	}
	response.JSON(c, http.StatusNoContent, nil)
}
