package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	tokens, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		case errors.Is(err, usecase.ErrAccountInactive):
			response.Error(c, http.StatusForbidden, "ACCOUNT_INACTIVE", "Account is inactive")
		default:
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Login failed")
		}
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
		"user": gin.H{
			"id":        tokens.User.ID,
			"name":      tokens.User.Name,
			"role":      tokens.User.Role,
			"school_id": tokens.User.SchoolID,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	tokens, err := h.uc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "TOKEN_INVALID", "Invalid or expired refresh token")
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
		"user": gin.H{
			"id":        tokens.User.ID,
			"name":      tokens.User.Name,
			"role":      tokens.User.Role,
			"school_id": tokens.User.SchoolID,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = c.ShouldBindJSON(&req)         //nolint:errcheck
	_ = h.uc.Logout(c.Request.Context(), req.RefreshToken) //nolint:errcheck
	c.Status(http.StatusNoContent)
}
