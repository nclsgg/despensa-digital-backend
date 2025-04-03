package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}
