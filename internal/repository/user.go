package repository

import (
	"context"
	"log"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/jmoiron/sqlx"
)

func (r *Repo) CreateUser(ctx context.Context, tx *sqlx.Tx, in *entity.User) (uint64, error) {
	var id uint64
	err := tx.QueryRowxContext(ctx, `insert into users(email, password, name) values ($1, $2, $3) returning id`, in.Email, in.Password, in.Name).Scan(&id)
	if err != nil {
		log.Println(ctx, err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindUser(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.dbRead.QueryRowxContext(ctx, `select id, email, password, name, created_at, updated_at, deleted_at from users where email = $1`, email).StructScan(&user)
	if err != nil {
		log.Println(ctx, err)
		return nil, err
	}

	return &user, nil
}
