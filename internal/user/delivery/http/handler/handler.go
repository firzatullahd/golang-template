package handler

import (
	"context"

	"github.com/firzatullahd/golang-template/internal/user/model"
)

//go:generate mockgen -source=handler.go -package=mock -destination=mock/mock.go
type Iservice interface {
	Register(ctx context.Context, in model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, in model.AuthRequest) (*model.AuthResponse, error)

	InitiateVerification(ctx context.Context, username string) error
	Verification(ctx context.Context, username, code string) error
}

type Handler struct {
	Service Iservice
}

func NewHandler(service Iservice) *Handler {
	return &Handler{
		Service: service,
	}
}
