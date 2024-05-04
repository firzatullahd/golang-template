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

func (r *Repo) CreateMatch(ctx context.Context, tx *sqlx.Tx, in *entity.Match) error {
	logCtx := fmt.Sprintf("%T.CreateMatch", r)
	_, err := tx.ExecContext(ctx, `insert into match(user_id, cat_id, match_user_id, match_cat_id, message) values ($1, $2, $3, $4, $5) returning id`, in.UserID, in.CatID, in.MatchUserID, in.MatchCatID, in.Message)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) DeleteMatch(ctx context.Context, tx *sqlx.Tx, id []uint64) error {
	logCtx := fmt.Sprintf("%T.CreateMatch", r)
	_, err := tx.ExecContext(ctx, `update match set deleted_at = now() where id in ($1)`, id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) UpdateMatch(ctx context.Context, tx *sqlx.Tx, in *model.InputUpdateMatch) error {
	logCtx := fmt.Sprintf("%T.UpdateMatch", r)

	query, args := buildQueryUpdateMatch(in)
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

func buildQueryUpdateMatch(in *model.InputUpdateMatch) (string, map[string]interface{}) {
	args := make(map[string]interface{}, 0)

	if in.Approval {
		args["approved"] = in.Approval
		args["rejected"] = !in.Approval
	} else {
		args["approved"] = !in.Approval
		args["rejected"] = in.Approval
	}

	// args["match_user_id"] = in.MatchUserId
	// args["match_cat_id"] = in.MatchCatId
	args["id"] = in.MatchId

	// query := "update match set is_approved = :approved, is_rejected = :rejected where match_user_id = :match_user_id and match_cat_id = :match_cat_id"
	query := "update match set is_approved = :approved, is_rejected = :rejected where id = :id"

	return query, args
}

func (r *Repo) FindMatch(ctx context.Context, filter *model.FilterFindMatch) ([]entity.Match, error) {
	logCtx := fmt.Sprintf("%T.FindMatch", r)

	var matches []entity.Match

	query, args := buildQueryFindMatch(filter)
	rows, err := r.dbRead.QueryxContext(ctx, query, args)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var match entity.Match
		if err := rows.StructScan(&match); err != nil {
			logger.Error(ctx, logCtx, err)
			return matches, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func buildQueryFindMatch(filter *model.FilterFindMatch) (string, map[string]interface{}) {

	args := make(map[string]interface{}, 0)
	var params []string
	if filter.CatId != nil {
		params = append(params, `(cat_id in (:cat_id) or match_cat_id in (:cat_id))`)
		args["cat_id"] = filter.CatId
	}

	if filter.UserId != nil {
		params = append(params, `(user_id = :user_id or match_user_id = :user_id)`)
		args["user_id"] = filter.UserId
	}

	if len(filter.ID) > 0 {
		params = append(params, `id in (:id)`)
		args["id"] = filter.ID
	}

	if filter.Approval != nil {
		if *filter.Approval {
			params = append(params, `is_approved = :approved`)
			args["approved"] = filter.Approval
		} else {
			params = append(params, `is_rejected = :rejected`)
			args["rejected"] = filter.Approval
		}
	} else if filter.PendingApproval {
		params = append(params, `is_approved = :false and is_rejected = :false`)
		args["false"] = false
	}

	query := fmt.Sprintf(`select id, user_id, cat_id, match_user_id, match_cat_id, is_approved, is_rejected, message, created_at, updated_at, deleted_at from match where deleted_at isnull and %s`, strings.Join(params, "and"))

	return query, args
}
