package model

import "github.com/golang-jwt/jwt/v5"

var JWT_SIGNATURE_KEY = []byte("cats-social")

type MyClaim struct {
	UserData
	jwt.RegisteredClaims
}

type UserData struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}
