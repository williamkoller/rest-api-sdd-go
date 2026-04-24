package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type ClassHandler struct {
	uc *usecase.ClassUseCase
}

func NewClassHandler(uc *usecase.ClassUseCase) *ClassHandler {
	return &ClassHandler{uc: uc}
}

func (h *ClassHandler) List(c *gin.Context) {
	var year *int
	if v := c.Query("academic_year"); v != "" {
		y, _ := strconv.Atoi(v) //nolint:errcheck
		year = &y
	}
	var active *bool
	if v := c.Query("active"); v != "" {
		b, _ := strconv.ParseBool(v) //nolint:errcheck
		active = &b
	}
	classes, err := h.uc.FindByUnitID(c.Request.Context(), c.Param("unit_id"), year, active)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch classes")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapClassesWithCount(classes))
}

func (h *ClassHandler) Create(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		GradeLevel   string `json:"grade_level" binding:"required"`
		Shift        string `json:"shift" binding:"required"`
		AcademicYear int    `json:"academic_year" binding:"required"`
		ClassroomID  string `json:"classroom_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	class := &entity.Class{
		UnitID:       c.Param("unit_id"),
		ClassroomID:  req.ClassroomID,
		Name:         req.Name,
		GradeLevel:   req.GradeLevel,
		Shift:        entity.Shift(req.Shift),
		AcademicYear: req.AcademicYear,
	}
	if err := h.uc.Create(c.Request.Context(), class); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create class")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapClass(class))
}

func (h *ClassHandler) Get(c *gin.Context) {
	result, err := h.uc.FindByID(c.Request.Context(), c.Param("class_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Class not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch class")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapClassWithCount(result))
}

func (h *ClassHandler) Update(c *gin.Context) {
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	class, err := h.uc.Update(c.Request.Context(), c.Param("class_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Class not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update class")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapClass(class))
}
