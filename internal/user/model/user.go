package model

import "time"

var VerificationCounterPrefix = `verification:counter:%s`
var VerificationPrefix = `verification:%s`
var VerificationTTL = 30 * time.Minute
var VerificationMaxAttempt = 3

type RegisterRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	UserID      uint64 `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FilterFindUser struct {
	Username *string
	ID       []uint64
}

type EmailPayload struct {
	Email            string
	Name             string
	VerificationCode string
}
