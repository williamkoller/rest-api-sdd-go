package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type GradeHandler struct {
	uc *usecase.GradeUseCase
}

func NewGradeHandler(uc *usecase.GradeUseCase) *GradeHandler {
	return &GradeHandler{uc: uc}
}

func (h *GradeHandler) BatchUpsert(c *gin.Context) {
	var req struct {
		Subject string `json:"subject" binding:"required"`
		Period  string `json:"period" binding:"required"`
		Grades  []struct {
			StudentID string  `json:"student_id"`
			Value     float64 `json:"value"`
		} `json:"grades" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	grades := make([]repository.GradeInput, len(req.Grades))
	for i, g := range req.Grades {
		if g.Value < 0 || g.Value > 10 {
			response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "grade value must be between 0 and 10")
			return
		}
		grades[i] = repository.GradeInput{StudentID: g.StudentID, Value: g.Value}
	}

	teacherID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	count, err := h.uc.BatchUpsert(c.Request.Context(), c.Param("class_id"), req.Subject, req.Period, grades, teacherID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Not assigned to this class")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to record grades")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": count})
}

func (h *GradeHandler) GetByStudent(c *gin.Context) {
	filters := repository.GradeFilters{
		Subject: c.Query("subject"),
		Period:  c.Query("period"),
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	grades, err := h.uc.GetByStudent(c.Request.Context(), c.Param("student_id"), filters, userID, role)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Access denied")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch grades")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapGrades(grades))
}
