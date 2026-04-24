package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type ClassroomHandler struct {
	uc *usecase.ClassroomUseCase
}

func NewClassroomHandler(uc *usecase.ClassroomUseCase) *ClassroomHandler {
	return &ClassroomHandler{uc: uc}
}

func (h *ClassroomHandler) List(c *gin.Context) {
	classrooms, err := h.uc.FindByUnitID(c.Request.Context(), c.Param("unit_id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch classrooms")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapClassrooms(classrooms))
}

func (h *ClassroomHandler) Create(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Capacity int    `json:"capacity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	classroom := &entity.Classroom{
		UnitID:   c.Param("unit_id"),
		Code:     req.Code,
		Capacity: req.Capacity,
	}
	if err := h.uc.Create(c.Request.Context(), classroom); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create classroom")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapClassroom(classroom))
}

func (h *ClassroomHandler) Update(c *gin.Context) {
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	classroom, err := h.uc.Update(c.Request.Context(), c.Param("classroom_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Classroom not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update classroom")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapClassroom(classroom))
}
