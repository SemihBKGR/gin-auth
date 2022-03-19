package handle

import (
	"github.com/gin-gonic/gin"
	"time"
)

type ErrorResponse struct {
	Status    int    `json:"status,omitempty"`
	Path      string `json:"path,omitempty"`
	Method    string `json:"method,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func wrapError(err error, status int, c *gin.Context) *ErrorResponse {
	if err == nil {
		panic("error cannot be nil")
	}
	return &ErrorResponse{
		Status:    status,
		Path:      c.Request.RequestURI,
		Method:    c.Request.Method,
		Message:   err.Error(),
		Timestamp: time.Now().UnixMilli(),
	}
}

func wrapErrorAndSend(err error, status int, c *gin.Context) {
	c.JSON(status, wrapError(err, status, c))
}
