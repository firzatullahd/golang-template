package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
	"github.com/jmoiron/sqlx"
)

func (r *Repo) CreateUser(ctx context.Context, tx *sqlx.Tx, in *entity.User) (uint64, error) {
	logCtx := fmt.Sprintf("%T.CreateUser", r)
	var id uint64
	err := tx.QueryRowxContext(ctx, `insert into users(email, password, name) values ($1, $2, $3) returning id`, in.Email, in.Password, in.Name).Scan(&id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error) {
	logCtx := fmt.Sprintf("%T.FindUser", r)
	var users []entity.User

	query, args := buildQueryFindUser(in)

	rows, err := r.dbRead.QueryxContext(ctx, query, args)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.StructScan(&user)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func buildQueryFindUser(filter *model.FilterFindUser) (string, map[string]interface{}) {

	args := make(map[string]interface{}, 0)
	var params []string
	if filter.Email != nil {
		params = append(params, `email = :email`)
		args["email"] = filter.Email
	}

	if len(filter.ID) > 0 {
		params = append(params, `id in (:id)`)
		args["id"] = filter.ID
	}

	query := fmt.Sprintf(`select id, email, password, name, created_at, updated_at, deleted_at from users where deleted_at isnull and %s`, strings.Join(params, "and"))

	return query, args
}
