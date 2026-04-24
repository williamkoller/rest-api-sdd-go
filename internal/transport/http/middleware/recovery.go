package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", "error", r)
				response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
