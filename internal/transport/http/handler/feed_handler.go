package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type FeedHandler struct {
	uc *usecase.FeedUseCase
}

func NewFeedHandler(uc *usecase.FeedUseCase) *FeedHandler {
	return &FeedHandler{uc: uc}
}

func (h *FeedHandler) Create(c *gin.Context) {
	var req struct {
		UnitID   string `json:"unit_id"`
		Body     string `json:"body" binding:"required"`
		ImageURL string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	post := &entity.FeedPost{
		SchoolID: middleware.GetSchoolID(c),
		UnitID:   req.UnitID,
		AuthorID: middleware.GetUserID(c),
		Body:     req.Body,
		ImageURL: req.ImageURL,
	}
	result, err := h.uc.Publish(c.Request.Context(), post, middleware.GetRole(c))
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Only staff can publish posts")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to publish post")
		return
	}
	response.JSON(c, http.StatusCreated, dto.MapFeedPost(result))
}

func (h *FeedHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))         //nolint:errcheck
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20")) //nolint:errcheck
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	var unitID *string
	if v := c.Query("unit_id"); v != "" {
		unitID = &v
	}
	posts, total, err := h.uc.List(c.Request.Context(), c.Param("school_id"), unitID, page, perPage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch feed")
		return
	}
	response.Paginated(c, dto.MapFeedPosts(posts), page, perPage, total)
}

func (h *FeedHandler) Delete(c *gin.Context) {
	err := h.uc.Delete(c.Request.Context(), c.Param("feed_id"), middleware.GetUserID(c), middleware.GetRole(c))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Post not found")
		case errors.Is(err, usecase.ErrForbidden):
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Cannot delete this post")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete post")
		}
		return
	}
	response.JSON(c, http.StatusNoContent, nil)
}
