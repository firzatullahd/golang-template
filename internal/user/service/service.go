package service

import (
	"context"

	"github.com/firzatullahd/golang-template/config"
	"github.com/firzatullahd/golang-template/internal/user/entity"
	"github.com/firzatullahd/golang-template/internal/user/model"
	timeutils "github.com/firzatullahd/golang-template/utils/time"
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

type IEmailClient interface {
	SendEmail(ctx context.Context, input model.EmailPayload) error
}

type Service struct {
	conf        *config.Config
	repo        Irepository
	redisConn   *redis.Client
	time        timeutils.Time
	emailClient IEmailClient
}

func NewService(conf *config.Config, repo Irepository, redisConn *redis.Client, emailClient IEmailClient) Service {
	return Service{
		conf:        conf,
		repo:        repo,
		redisConn:   redisConn,
		time:        timeutils.Time{},
		emailClient: emailClient,
	}
}
