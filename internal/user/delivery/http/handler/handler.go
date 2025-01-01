package handler

import (
	"context"

	"github.com/firzatullahd/golang-template/internal/user/model"
)

type Iservice interface {
	Register(ctx context.Context, in model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, in model.AuthRequest) (*model.AuthResponse, error)
}

type Handler struct {
	Service Iservice
}

func NewHandler(service Iservice) *Handler {
	return &Handler{
		Service: service,
	}
}
