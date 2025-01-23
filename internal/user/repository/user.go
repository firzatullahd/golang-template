package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/firzatullahd/golang-template/internal/user/entity"
	"github.com/firzatullahd/golang-template/internal/user/model"
	customerror "github.com/firzatullahd/golang-template/internal/user/model/error"
	"github.com/jmoiron/sqlx"
)

func (r *Repo) CreateUser(ctx context.Context, tx *sqlx.Tx, in entity.User) (uint64, error) {
	var id uint64

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Insert("users").Columns("username", "password", "name", "state").Values(in.Username, in.Password, in.Name, entity.UserStateRegistered).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}
	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindUser(ctx context.Context, filter *model.FilterFindUser) (*entity.User, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryfind := sq.Select("id", "username", "password", "name", "state", "id_card_no", "id_card_file", "created_at", "updated_at").From("users")

	switch {
	case filter.Username != nil:
		queryfind = queryfind.Where(squirrel.Eq{"username": *filter.Username})
	case filter.ID != nil:
		queryfind = queryfind.Where(squirrel.Eq{"id": filter.ID})
	}

	query, args, err := queryfind.ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User

	if err := r.dbRead.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) UpdateUser(ctx context.Context, tx *sqlx.Tx, filter *model.FilterFindUser, in map[string]any) error {
	in["updated_at"] = time.Now()
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := sq.Update("users").SetMap(in)
	switch {
	case filter.Username != nil:
		q = q.Where(squirrel.Eq{"username": *filter.Username})
	case filter.ID != nil:
		q = q.Where(squirrel.Eq{"id": filter.ID})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return customerror.ErrNoResourceUpdated
	}

	return nil
}
