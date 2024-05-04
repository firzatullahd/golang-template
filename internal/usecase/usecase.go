package usecase

import (
	"context"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type Irepository interface {
	WithTransaction() (*sqlx.Tx, error)

	CreateUser(ctx context.Context, tx *sqlx.Tx, in *entity.User) (uint64, error)
	FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error)

	CreateCat(ctx context.Context, tx *sqlx.Tx, in *entity.Cat) (uint64, error)
	FindCat(ctx context.Context, filter *model.FilterFindCat) ([]entity.Cat, error)
	DeleteCat(ctx context.Context, tx *sqlx.Tx, catId, userId uint64) error
	UpdateCat(ctx context.Context, tx *sqlx.Tx, in *model.InputUpdateCat) error

	CreateMatch(ctx context.Context, tx *sqlx.Tx, in *entity.Match) error
	UpdateMatch(ctx context.Context, tx *sqlx.Tx, in *model.InputUpdateMatch) error
	DeleteMatch(ctx context.Context, tx *sqlx.Tx, id []uint64) error
	FindMatch(ctx context.Context, filter *model.FilterFindMatch) ([]entity.Match, error)
}

type Usecase struct {
	conf *config.Config
	repo Irepository
}

func NewUsecase(conf *config.Config, repo Irepository) *Usecase {
	return &Usecase{
		conf: conf,
		repo: repo,
	}
}
