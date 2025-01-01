package model

import "github.com/golang-jwt/jwt/v5"

type MyClaim struct {
	UserData
	jwt.RegisteredClaims
}

type UserData struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
}
