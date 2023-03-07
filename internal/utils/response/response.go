package response

import (
	"fmt"
	"strings"

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
		case "alphanum":
			message = fmt.Sprintf("%s must be alphanumeric", validator.Field())
		case "gte":
			message = fmt.Sprintf("%s must be greater or equal than %s", validator.Field(), validator.Param())
		case "datetime":
			switch validator.Param() {
			case "01/06":
				message = fmt.Sprintf("%s must be MM/YY", validator.Field())
			case "2006-01-02":
				message = fmt.Sprintf("%s must be YYYY-MM-DD", validator.Field())
			}
		case "url":
			message = fmt.Sprintf("%s must be a URL", validator.Field())
		case "oneof":
			vals := validator.Param()
			valueList := strings.Split(vals, " ")
			valMessage := ""
			for i, v := range valueList {
				if i == len(valueList)-1 {
					valMessage += fmt.Sprintf("or %s", v)
					continue
				}

				valMessage += fmt.Sprintf("%s, ", v)
			}
			message = fmt.Sprintf("%s must be either %s", validator.Field(), valMessage)
		case "required_without":
			message = fmt.Sprintf("%s is required", validator.Field())
		case "required_with":
			message = fmt.Sprintf("%s is required", validator.Field())
		}
	}

	c.JSON(statusCode, Response{
		Code:    "BAD_REQUEST",
		Message: message,
	})
}
