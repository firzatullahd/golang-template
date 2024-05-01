package repository

import (
	"context"
	"log"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/jmoiron/sqlx"
)

func (r *Repo) CreateUser(ctx context.Context, tx *sqlx.Tx, in *entity.User) (int, error) {
	var id int
	err := tx.QueryRowxContext(ctx, `insert into users(email, password, name) values ($1, $2, $3) returning id`, in.Email, in.Password, in.Name).Scan(&id)
	if err != nil {
		log.Println(ctx, err)
		return 0, err
	}

	return id, nil
}
