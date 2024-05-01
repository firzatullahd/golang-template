package usecase

import (
	"context"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Irepository interface {
	WithTransaction() (*sqlx.Tx, error)
	CreateUser(ctx context.Context, tx *sqlx.Tx, in *entity.User) (int, error)
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
