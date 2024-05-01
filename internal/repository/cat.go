package repository

import (
	"context"
	"fmt"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
	"github.com/jmoiron/sqlx"
)

func (r *Repo) CreateCat(ctx context.Context, tx *sqlx.Tx, in *entity.Cat) (uint64, error) {
	logCtx := fmt.Sprintf("%T.CreateCat", r)
	var id uint64
	err := tx.QueryRowxContext(ctx, `insert into cat(user_id, name, sex, race, image_urls, age, description) values ($1, $2, $3, $4, $5, $6, $7) returning id`, in.UserID, in.Name, in.Sex, in.Race, in.ImageUrls).Scan(&id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindCat(ctx context.Context, filter model.FilterFindCat) ([]entity.Cat, error) {
	logCtx := fmt.Sprintf("%T.FindCat", r)
	// query := `select id, user_id, name, sex, race, image_urls, age, description, has_matched created_at, updated_at, deleted_at from cat`

	query, args := buildQueryFindCat(filter)

	var cats []entity.Cat
	rows, err := r.dbRead.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat entity.Cat
		if err := rows.StructScan(&cat); err != nil {
			logger.Error(ctx, logCtx, err)
			return cats, err
		}

		cats = append(cats, cat)
	}

	return cats, nil
}
func buildQueryFindCat(filter model.FilterFindCat) (string, []interface{}) {
	query := `select id, user_id, name, sex, race, image_urls, age, description, has_matched created_at, updated_at, deleted_at from cat where 1=1`

	if filter.Limit == 0 {
		filter.Limit = 5
	}

	args := []interface{}{}

	if filter.ID > 0 {
		query += `and id = ?`
		args = append(args, filter.ID)
	}

	if filter.Sex != "" {
		query += `and sex = ?`
		args = append(args, filter.Sex)
	}

	if filter.Race != "" {
		query += `and race = ?`
		args = append(args, filter.Race)
	}

	if filter.HasMatched != nil {
		query += `and has_matched = ?`
		args = append(args, filter.HasMatched)
	}

	if filter.Age > 0 {
		query += `and age = ?`
		args = append(args, filter.Age)
	}

	if filter.Owned {
		query += `and user_id = ?`
		args = append(args, filter.UserID)
	}

	if filter.SearchName != "" {
		query += `and name = ilike '%?%'`
		args = append(args, filter.SearchName)
	}

	return query, args
}

func (r *Repo) DeleteCat(ctx context.Context, tx *sqlx.Tx, id uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteCat", r)
	_, err := tx.ExecContext(ctx, `update cat set deleted_at = now() where id = $1`, id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) UpdateCat(ctx context.Context, tx *sqlx.Tx, in *entity.Cat) error {
	logCtx := fmt.Sprintf("%T.UpdateCat", r)
	query := `update cat set name = $2, sex = $3, race = $4, image_urls = $5, updated_at = now() where id = $1`
	_, err := tx.ExecContext(ctx, query, in.ID, in.Name, in.Sex, in.Race, in.ImageUrls)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}
