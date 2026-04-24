package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type AttendanceHandler struct {
	uc *usecase.AttendanceUseCase
}

func NewAttendanceHandler(uc *usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{uc: uc}
}

func (h *AttendanceHandler) BatchRecord(c *gin.Context) {
	var req struct {
		Date    string `json:"date" binding:"required"`
		Records []struct {
			StudentID string `json:"student_id"`
			Status    string `json:"status"`
			Note      string `json:"note"`
		} `json:"records" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid date format, use YYYY-MM-DD")
		return
	}

	records := make([]repository.AttendanceInput, len(req.Records))
	for i, r := range req.Records {
		records[i] = repository.AttendanceInput{
			StudentID: r.StudentID,
			Status:    entity.AttendanceStatus(r.Status),
			Note:      r.Note,
		}
	}

	teacherID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	count, err := h.uc.BatchRecord(c.Request.Context(), c.Param("class_id"), date, records, teacherID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Not assigned to this class")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to record attendance")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"recorded": count})
}

func (h *AttendanceHandler) GetByStudent(c *gin.Context) {
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

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	summary, err := h.uc.GetByStudent(c.Request.Context(), c.Param("student_id"), from, to, userID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Access denied")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch attendance")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapAttendanceSummary(summary))
}
