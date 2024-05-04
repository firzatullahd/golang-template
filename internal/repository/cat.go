package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
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

func (r *Repo) FindCat(ctx context.Context, filter *model.FilterFindCat) ([]entity.Cat, error) {
	logCtx := fmt.Sprintf("%T.FindCat", r)
	query, args := buildQueryFindCat(filter)

	var cats []entity.Cat
	rows, err := r.dbRead.QueryxContext(ctx, query, args)
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

func buildQueryFindCat(filter *model.FilterFindCat) (string, map[string]interface{}) {
	query := `select id, user_id, name, sex, race, image_urls, age, description, has_matched created_at, updated_at, deleted_at from cat where deleted_at isnull `

	if filter.Limit == 0 {
		filter.Limit = 5
	}

	args := make(map[string]interface{}, 0)

	if filter.ID != nil {
		query += `and id in (:id)`
		args["id"] = filter.ID
	}

	if filter.Sex != nil {
		query += `and sex = :sex`
		args["sex"] = filter.Sex
	}

	if filter.Race != nil {
		query += `and race = :race`
		args["race"] = filter.Race
	}

	if filter.HasMatched != nil {
		query += `and has_matched = :has_matched`
		args["has_matched"] = filter.HasMatched
	}

	if filter.Age != nil {
		switch *filter.Age {
		case ">4":
			query += `and age >= :age`
		case "<4":
			query += `and age <= :age`
		case "4":
			query += `and age = :age`
		}

		args["age"] = filter.Age
	}

	if filter.UserID != nil {
		query += `and user_id = ?`
		args["user_id"] = filter.UserID
	}

	if filter.SearchName != nil {
		query += `and name = ilike '%:search_name%'`
		args["search_name"] = filter.SearchName
	}

	query += fmt.Sprintf(`limit %d offset %d`, filter.Limit, filter.Offset)

	return query, args
}

func (r *Repo) DeleteCat(ctx context.Context, tx *sqlx.Tx, catId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteCat", r)
	res, err := tx.ExecContext(ctx, `update cat set deleted_at = now() where id = $1 and user_id = $2`, catId, userId)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	if affected <= 0 {
		return error_envelope.ErrNotFound
	}

	return nil
}

func (r *Repo) UpdateCat(ctx context.Context, tx *sqlx.Tx, in *model.InputUpdateCat) error {
	logCtx := fmt.Sprintf("%T.UpdateCat", r)
	query, args := buildQueryUpdateCat(in)
	res, err := tx.ExecContext(ctx, query, args)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	if affected <= 0 {
		return error_envelope.ErrNotFound
	}

	return nil
}

func buildQueryUpdateCat(in *model.InputUpdateCat) (string, map[string]interface{}) {
	args := make(map[string]interface{}, 0)
	var params []string
	if in.Name != nil {
		params = append(params, `name = :name`)
		args["name"] = in.Name
	}

	if in.Sex != nil {
		params = append(params, `sex = :sex`)
		args["sex"] = in.Sex
	}

	if in.Race != nil {
		params = append(params, `race = :race`)
		args["race"] = in.Race
	}

	if in.ImageUrls != nil {
		params = append(params, `image_urls = :image_urls`)
		args["image_urls"] = in.ImageUrls
	}

	if in.Age != nil {
		params = append(params, `age = :age`)
		args["age"] = in.Age
	}

	if in.Description != nil {
		params = append(params, `description = :desc`)
		args["desc"] = in.Description
	}

	args["id"] = in.ID
	args["user_id"] = in.UserID

	query := fmt.Sprintf("update cat set %s where id = :id and user_id = :user_id", strings.Join(params, ","))

	return query, args
}
