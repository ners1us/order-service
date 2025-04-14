package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
