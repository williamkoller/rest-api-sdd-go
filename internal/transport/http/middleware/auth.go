package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/response"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeySchoolID contextKey = "school_id"
	ContextKeyRole     contextKey = "role"
)

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "TOKEN_INVALID", "Missing or malformed authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			code := "TOKEN_INVALID"
			if strings.Contains(err.Error(), "expired") {
				code = "TOKEN_EXPIRED"
			}
			response.Error(c, http.StatusUnauthorized, code, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "TOKEN_INVALID", "Invalid token claims")
			c.Abort()
			return
		}

		c.Set(string(ContextKeyUserID), claims["user_id"])
		c.Set(string(ContextKeySchoolID), claims["school_id"])
		c.Set(string(ContextKeyRole), claims["role"])
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(string(ContextKeyUserID))
	s, _ := v.(string)
	return s
}

func GetSchoolID(c *gin.Context) string {
	v, _ := c.Get(string(ContextKeySchoolID))
	s, _ := v.(string)
	return s
}

func GetRole(c *gin.Context) string {
	v, _ := c.Get(string(ContextKeyRole))
	s, _ := v.(string)
	return s
}
