package middleware

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const emptyToken = ""

func JWTAuthorization(c *gin.Context) {
	auth := c.GetHeader("authorization")

	if auth == emptyToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: errs.ErrInvalidToken.Error(),
		})
		return
	}

	auth = strings.Replace(auth, "Bearer ", "", -1)

	parsedToken, err := jwttoken.ValidateToken(auth, config.SecretKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	c.Set("userId", parsedToken.UserId)
}

func JWTValidateRefreshToken(c *gin.Context) {
	auth := c.GetHeader("authorization")

	if auth == emptyToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: errs.ErrInvalidToken.Error(),
		})
		return
	}

	auth = strings.Replace(auth, "Bearer ", "", -1)

	parsedToken, err := jwttoken.ValidateRefreshToken(auth, config.SecretKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	c.Set("userId", parsedToken.UserId)
}
