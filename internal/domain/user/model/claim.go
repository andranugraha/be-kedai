package model

import "github.com/golang-jwt/jwt/v4"

type Claim struct {
	UserId int `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}