package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, statusCode int, code, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}
