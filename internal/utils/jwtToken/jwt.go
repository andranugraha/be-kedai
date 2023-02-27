package jwttoken

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken(user *model.User) (string, error) {

	accessTime := ParseTokenAgeFromENV(config.GetEnv("ACCESS_TOKEN_AGE", ""), "access")
	claims := &model.Claim{
		UserId:    user.ID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	return tokenString, nil
}

func GenerateRefreshToken(user *model.User) (string, error) {
	refreshTime := ParseTokenAgeFromENV(config.GetEnv("REFRESH_TOKEN_AGE", ""), "refresh")
	claims := &model.Claim{
		UserId:    user.ID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Kedai",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTime)),
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

func ParseTokenAgeFromENV(age string, tokenType string) time.Duration {
	ageNum, _ := strconv.Atoi(age)

	if tokenType == "access" {
		return time.Minute * time.Duration(ageNum)
	}

	return time.Hour * time.Duration(ageNum)
}
