package model

import "github.com/golang-jwt/jwt/v4"

type Claim struct {
<<<<<<< HEAD
	UserId int `json:"userId"`
=======
	UserId    int    `json:"userId"`
	TokenType string `json:"tokenType"`
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
	jwt.RegisteredClaims
}
