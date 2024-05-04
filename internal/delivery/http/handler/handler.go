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
	FindCat(ctx context.Context, in *model.FindCatRequest) ([]model.FindCatResponse, error)
	UpdateCat(ctx context.Context, in *model.UpdateCatRequest) error

	// Manage Match
	CreateMatch(ctx context.Context, in *model.CreateMatchRequest) error
	FindMatch(ctx context.Context, userId uint64) ([]model.FindMatchResponse, error)
	DeleteMatch(ctx context.Context, matchId, userId uint64) error
	RejectMatch(ctx context.Context, matchId, userId uint64) error
	ApproveMatch(ctx context.Context, matchId, userId uint64) error
}

type Handler struct {
	Usecase IUsecase
}

func NewHandler(usecase IUsecase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}
