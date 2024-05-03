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
	FindUser(ctx context.Context, email string) (*entity.User, error)

	CreateCat(ctx context.Context, tx *sqlx.Tx, in *entity.Cat) (uint64, error)
	FindCat(ctx context.Context, filter *model.FilterFindCat) ([]entity.Cat, error)
	DeleteCat(ctx context.Context, tx *sqlx.Tx, catId, userId uint64) error
	UpdateCat(ctx context.Context, tx *sqlx.Tx, in *entity.Cat) error
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
