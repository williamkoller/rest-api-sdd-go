package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Data  any         `json:"data"`
	Error *apiError   `json:"error,omitempty"`
	Meta  *pagination `json:"meta,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, envelope{Data: data})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, envelope{Data: nil, Error: &apiError{Code: code, Message: message}})
}

func Paginated(c *gin.Context, data any, page, perPage int, total int64) {
	c.JSON(http.StatusOK, envelope{
		Data: data,
		Meta: &pagination{Page: page, PerPage: perPage, Total: total},
	})
}
