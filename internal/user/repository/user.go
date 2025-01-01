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
	// logCtx := fmt.Sprintf("%T.CreateUser", r)
	var id uint64

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Insert("users").Columns("username", "password", "name", "state").Values(in.Username, in.Password, in.Name, entity.UserStateRegistered).ToSql()
	if err != nil {
		return 0, err
	}
	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error) {
	// logCtx := fmt.Sprintf("%T.FindUser", r)
	var users []entity.User

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Select("id", "username", "password", "name", "state", "id_card_no", "id_card_file", "created_at", "updated_at").From("users").ToSql()
	if err != nil {
		return users, err
	}

	rows, err := r.dbRead.QueryxContext(ctx, query, args)
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.StructScan(&user)
		if err != nil {
			// logger.Error(ctx, logCtx, err)
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) FindUser(ctx context.Context, in *model.FilterFindUser) (*entity.User, error) {
	// logCtx := fmt.Sprintf("%T.FindUser", r)

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryfind := sq.Select("id", "username", "password", "name", "state", "id_card_no", "id_card_file", "created_at", "updated_at").From("users")

	if in.Username != nil {
		queryfind = queryfind.Where(squirrel.Eq{"username": in.Username})
	}

	switch {
	case in.Username != nil:
		queryfind = queryfind.Where(squirrel.Eq{"username": in.Username})
	case len(in.ID) > 0:
		queryfind = queryfind.Where(squirrel.Eq{"id": in.ID})
	}

	query, args, err := queryfind.ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User

	if err := r.dbRead.QueryRowxContext(ctx, query, args).StructScan(&user); err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &user, nil
}

func (r *Repo) UpdateUser(ctx context.Context, tx *sqlx.Tx, userID uint64, in map[string]any) error {
	// logCtx := fmt.Sprintf("%T.UpdateUser", r)

	in["updated_at"] = time.Now()

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.Update("users").SetMap(in).Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		// logger.Error(ctx, logCtx, err)
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
