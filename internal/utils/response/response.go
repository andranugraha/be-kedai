package response

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func ErrorValidator(c *gin.Context, statusCode int, err error) {
	var message string
	castedErr, _ := err.(validator.ValidationErrors)
	for _, validator := range castedErr {
		switch validator.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", validator.Field())
		case "email":
			message = fmt.Sprintf("%s must be an email format", validator.Field())
		case "min":
			message = fmt.Sprintf("%s must be greater than %s", validator.Field(), validator.Param())
		case "max":
			message = fmt.Sprintf("%s must be shorter than %s", validator.Field(), validator.Param())
		case "len":
			message = fmt.Sprintf("%s must be %s characters", validator.Field(), validator.Param())
		case "numeric":
			message = fmt.Sprintf("%s must be numeric", validator.Field())
		}
	}

	c.JSON(statusCode, Response{
		Code:    "BAD_REQUEST",
		Message: message,
	})
}
