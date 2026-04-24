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

type UnitHandler struct {
	uc *usecase.UnitUseCase
}

func NewUnitHandler(uc *usecase.UnitUseCase) *UnitHandler {
	return &UnitHandler{uc: uc}
}

func (h *UnitHandler) List(c *gin.Context) {
	var active *bool
	if v := c.Query("active"); v != "" {
		b, _ := strconv.ParseBool(v) //nolint:errcheck
		active = &b
	}
	units, err := h.uc.FindBySchoolID(c.Request.Context(), c.Param("school_id"), active)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch units")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapUnits(units))
}

func (h *UnitHandler) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
		City    string `json:"city" binding:"required"`
		State   string `json:"state" binding:"required,len=2"`
		ZipCode string `json:"zip_code" binding:"required"`
		Phone   string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	unit := &entity.Unit{
		SchoolID: c.Param("school_id"),
		Name:     req.Name,
		Address:  req.Address,
		City:     req.City,
		State:    req.State,
		ZipCode:  req.ZipCode,
		Phone:    req.Phone,
	}
	if err := h.uc.Create(c.Request.Context(), unit); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create unit")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapUnit(unit))
}

func (h *UnitHandler) Get(c *gin.Context) {
	unit, err := h.uc.FindByID(c.Request.Context(), c.Param("unit_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Unit not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch unit")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapUnit(unit))
}

func (h *UnitHandler) Update(c *gin.Context) {
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	unit, err := h.uc.Update(c.Request.Context(), c.Param("unit_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Unit not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update unit")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapUnit(unit))
}
