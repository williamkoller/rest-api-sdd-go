package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		schoolID := GetSchoolID(c)
		role := GetRole(c)
		// super_admin has no school_id in token — allowed to pass without tenant scoping
		if schoolID == "" && role != "super_admin" {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "No school context")
			c.Abort()
			return
		}
		c.Next()
	}
}
