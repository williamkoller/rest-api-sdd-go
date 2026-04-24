package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type MenuHandler struct {
	uc *usecase.MenuUseCase
}

func NewMenuHandler(uc *usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{uc: uc}
}

type menuItemRequest struct {
	DayOfWeek   string `json:"day_of_week" binding:"required"`
	MealType    string `json:"meal_type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req struct {
		WeekStart string            `json:"week_start" binding:"required"`
		Items     []menuItemRequest `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	weekStart, err := time.Parse("2006-01-02", req.WeekStart)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid week_start format")
		return
	}
	items := make([]*entity.MenuItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &entity.MenuItem{
			DayOfWeek:   entity.DayOfWeek(item.DayOfWeek),
			MealType:    entity.MealType(item.MealType),
			Description: item.Description,
		}
	}
	menu, err := h.uc.Publish(c.Request.Context(), c.Param("unit_id"), weekStart, items, middleware.GetUserID(c), middleware.GetRole(c))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Only staff can publish menus")
		case errors.Is(err, usecase.ErrMenuAlreadyExists):
			response.Error(c, http.StatusConflict, "MENU_ALREADY_EXISTS", "A menu already exists for this week")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to publish menu")
		}
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapMenu(menu))
}

func (h *MenuHandler) Get(c *gin.Context) {
	weekStart, err := time.Parse("2006-01-02", c.Query("week_start"))
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "week_start query param required (YYYY-MM-DD)")
		return
	}
	menu, err := h.uc.GetMenu(c.Request.Context(), c.Param("unit_id"), weekStart)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Menu not found for this week")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch menu")
		return
	}
	response.JSON(c, http.StatusOK, dto.MapMenu(menu))
}

func (h *MenuHandler) Update(c *gin.Context) {
	var req struct {
		Items []menuItemRequest `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	items := make([]*entity.MenuItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &entity.MenuItem{
			DayOfWeek:   entity.DayOfWeek(item.DayOfWeek),
			MealType:    entity.MealType(item.MealType),
			Description: item.Description,
		}
	}
	if err := h.uc.ReplaceMenu(c.Request.Context(), c.Param("menu_id"), items); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Menu not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to replace menu items")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"updated": true})
}
