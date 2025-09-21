package main

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}
