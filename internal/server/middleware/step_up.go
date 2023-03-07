package middleware

import (
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StepUp(c *gin.Context) {
	var (
		level          = c.GetInt("level")
		steppedUpLevel = 1
	)

	if level < steppedUpLevel {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: errs.ErrUnauthorized.Error(),
		})
		return
	}
}
