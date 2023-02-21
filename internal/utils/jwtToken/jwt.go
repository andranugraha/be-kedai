package jwttoken

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken(user *model.User) (*dto.Token, error) {
	claims := &model.Claim{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	result := &dto.Token{
		AccessToken: tokenString,
	}

	return result, nil
}

func ValidateToken(token string, secretKey string) (*model.Claim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &model.Claim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, jwt.ErrTokenMalformed
			}
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, jwt.ErrTokenExpired
			}
		}
	}

	return parsedToken.Claims.(*model.Claim), nil
}
