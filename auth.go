package main

import (
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(EnvOrDefault("JWT_SECRET", ""))
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}
