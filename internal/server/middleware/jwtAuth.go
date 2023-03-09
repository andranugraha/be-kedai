package middleware

import (
	"errors"
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    code.TOKEN_EXPIRED,
				Message: errs.ErrExpiredToken.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	c.Set("userId", parsedToken.UserId)
	c.Set("level", parsedToken.Level)
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
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    code.TOKEN_EXPIRED,
				Message: errs.ErrExpiredToken.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	c.Set("userId", parsedToken.UserId)
}
