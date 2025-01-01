package service

import (
	"context"

	"github.com/firzatullahd/golang-template/config"
	"github.com/firzatullahd/golang-template/internal/user/entity"
	"github.com/firzatullahd/golang-template/internal/user/model"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Irepository interface {
	WithTransaction() (*sqlx.Tx, error)

	CreateUser(ctx context.Context, tx *sqlx.Tx, in entity.User) (uint64, error)
	FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error)
	FindUser(ctx context.Context, in *model.FilterFindUser) (*entity.User, error)
	UpdateUser(ctx context.Context, tx *sqlx.Tx, userID uint64, in map[string]any) error
}

type Service struct {
	conf      *config.Config
	repo      Irepository
	redisConn *redis.Client
}

func NewService(conf *config.Config, repo Irepository, redisConn *redis.Client) Service {
	return Service{
		conf:      conf,
		repo:      repo,
		redisConn: redisConn,
	}
}
