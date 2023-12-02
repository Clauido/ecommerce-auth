package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId int32 `json:"userId"`
	jwt.StandardClaims
}