package handler

import (
	"context"

	"github.com/firzatullahd/cats-social-api/internal/model"
)

type IUsecase interface {
	// Auth
	Register(ctx context.Context, in *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error)

	// Manage Cat
	CreateCat(ctx context.Context, in *model.CreateCatRequest, userId uint64) (*model.CreateCatResponse, error)
	DeleteCat(ctx context.Context, catId, userId uint64) error
	FindCat(ctx context.Context, in *model.FilterFindCat) error
}

type Handler struct {
	Usecase IUsecase
}

func NewHandler(usecase IUsecase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}
