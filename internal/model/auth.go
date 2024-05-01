package model

import "github.com/golang-jwt/jwt/v5"

var JWT_SIGNATURE_KEY = []byte("cats-social")

type MyClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
