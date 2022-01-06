package api

import (
	"github.com/gin-gonic/gin"
)

// errorResponse represent entities of error with response
type errorResponse struct {
	Error string `json:"error"` // Error information
}

// newErrorResponse creates a new response with error
func newErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Error: msg,
	})
}
