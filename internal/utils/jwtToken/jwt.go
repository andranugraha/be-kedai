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
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "Kedai",
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