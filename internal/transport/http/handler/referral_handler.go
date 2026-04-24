package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/dto"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type ReferralHandler struct {
	uc *usecase.ReferralUseCase
}

func NewReferralHandler(uc *usecase.ReferralUseCase) *ReferralHandler {
	return &ReferralHandler{uc: uc}
}

func (h *ReferralHandler) GetMyLink(c *gin.Context) {
	link, err := h.uc.GetMyLink(c.Request.Context(), middleware.GetUserID(c), middleware.GetSchoolID(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get referral link")
		return
	}
	response.JSON(c, http.StatusOK, link)
}

func (h *ReferralHandler) ListReferrals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))         //nolint:errcheck
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20")) //nolint:errcheck
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	referrals, total, err := h.uc.GetReport(
		c.Request.Context(),
		c.Param("school_id"),
		c.Query("status"),
		middleware.GetRole(c),
		page,
		perPage,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrForbidden) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Only school admin can view referral reports")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch referrals")
		return
	}
	response.Paginated(c, dto.MapReferrals(referrals), page, perPage, total)
}
