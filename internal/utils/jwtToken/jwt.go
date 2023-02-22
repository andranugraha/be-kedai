package jwttoken

import (
	"kedai/backend/be-kedai/config"
<<<<<<< HEAD
	"kedai/backend/be-kedai/internal/domain/user/dto"
=======
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

<<<<<<< HEAD
func GenerateAccessToken(user *model.User) (*dto.Token, error) {
	claims := &model.Claim{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
=======
func GenerateAccessToken(user *model.User) (string, error) {
	claims := &model.Claim{
		UserId:    user.ID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

<<<<<<< HEAD
	result := &dto.Token{
		AccessToken: tokenString,
	}

	return result, nil
=======
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
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
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

<<<<<<< HEAD
	return parsedToken.Claims.(*model.Claim), nil
=======
	parsedClaim := parsedToken.Claims.(*model.Claim)
	if parsedClaim.TokenType != "access" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return parsedClaim, nil
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
}
