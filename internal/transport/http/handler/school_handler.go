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

type SchoolHandler struct {
	uc *usecase.SchoolUseCase
}

func NewSchoolHandler(uc *usecase.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{uc: uc}
}

func (h *SchoolHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))      //nolint:errcheck
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20")) //nolint:errcheck
	var active *bool
	if v := c.Query("active"); v != "" {
		b, _ := strconv.ParseBool(v) //nolint:errcheck
		active = &b
	}

	schools, total, err := h.uc.FindAll(c.Request.Context(), page, perPage, active)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch schools")
		return
	}
	response.Paginated(c, dto.MapSchools(schools), page, perPage, total)
}

func (h *SchoolHandler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		CNPJ  string `json:"cnpj" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	school := &entity.School{
		Name:  req.Name,
		CNPJ:  req.CNPJ,
		Email: req.Email,
		Phone: req.Phone,
	}
	if err := h.uc.Create(c.Request.Context(), school); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create school")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapSchool(school))
}

func (h *SchoolHandler) Get(c *gin.Context) {
	school, err := h.uc.FindByID(c.Request.Context(), c.Param("school_id"))
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "School not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch school")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapSchool(school))
}

func (h *SchoolHandler) Update(c *gin.Context) {
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	school, err := h.uc.Update(c.Request.Context(), c.Param("school_id"), updates)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "School not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update school")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapSchool(school))
}
