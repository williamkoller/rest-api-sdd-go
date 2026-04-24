package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/cache"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewHealthHandler(db *gorm.DB, cache cache.Cache) *HealthHandler {
	return &HealthHandler{db: db, cache: cache}
}

func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "ok"
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "error"
	}

	cacheStatus := "ok"
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := h.cache.Set(ctx, "health:ping", []byte("1"), 5*time.Second); err != nil {
		cacheStatus = "error"
	}

	status := http.StatusOK
	if dbStatus != "ok" || cacheStatus != "ok" {
		status = http.StatusServiceUnavailable
	}

	response.JSON(c, status, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"db":      dbStatus,
		"cache":   cacheStatus,
	})
}
