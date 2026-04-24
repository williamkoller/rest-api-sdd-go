package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type StudentHandler struct {
	uc *usecase.StudentUseCase
}

func NewStudentHandler(uc *usecase.StudentUseCase) *StudentHandler {
	return &StudentHandler{uc: uc}
}

func (h *StudentHandler) List(c *gin.Context) {
	students, err := h.uc.FindByClassID(c.Request.Context(), c.Param("class_id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch students")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapStudents(students))
}

func (h *StudentHandler) Enroll(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)

	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	req := usecase.EnrollRequest{EnrolledAt: time.Now()}
	if v, ok := body["student_id"].(string); ok {
		req.StudentID = v
	}
	if v, ok := body["name"].(string); ok {
		req.Name = v
	}
	if v, ok := body["cpf"].(string); ok {
		req.CPF = v
	}
	if v, ok := body["birth_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			req.BirthDate = t
		}
	}
	if v, ok := body["academic_year"].(float64); ok {
		req.AcademicYear = int(v)
	}
	if v, ok := body["enrolled_at"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			req.EnrolledAt = t
		}
	}

	if req.StudentID == "" && req.Name == "" {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "student_id or name is required")
		return
	}

	enrollment, err := h.uc.Enroll(c.Request.Context(), c.Param("class_id"), schoolID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Student not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to enroll student")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapEnrollment(enrollment))
}

func (h *StudentHandler) Get(c *gin.Context) {
	student, err := h.uc.FindByID(c.Request.Context(), c.Param("student_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Student not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch student")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapStudent(student))
}

func (h *StudentHandler) Update(c *gin.Context) {
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	student, err := h.uc.Update(c.Request.Context(), c.Param("student_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Student not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update student")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapStudent(student))
}

func (h *StudentHandler) Unenroll(c *gin.Context) {
	err := h.uc.Unenroll(c.Request.Context(), c.Param("class_id"), c.Param("student_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Enrollment not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to unenroll student")
		return
	}
	c.Status(http.StatusNoContent)
}
