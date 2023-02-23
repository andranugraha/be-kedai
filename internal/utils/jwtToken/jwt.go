package jwttoken

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken(user *model.User) (string, error) {
	claims := &model.Claim{
		UserId:    user.ID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	return tokenString, nil
}

func GenerateRefreshToken(user *model.User) (string, error) {
	claims := &model.Claim{
		UserId:    user.ID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	return tokenString, nil
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

	parsedClaim := parsedToken.Claims.(*model.Claim)
	if parsedClaim.TokenType != "access" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return parsedClaim, nil
}
