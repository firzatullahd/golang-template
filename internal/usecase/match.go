package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/firzatullahd/cats-social-api/internal/utils/constant"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
	"golang.org/x/sync/errgroup"
)

func (u *Usecase) CreateMatch(ctx context.Context, in *model.CreateMatchRequest) error {
	logCtx := fmt.Sprintf("%T.CreateMatch", u)

	cats, err := u.repo.FindCat(ctx, &model.FilterFindCat{
		Limit: 1,
		ID:    []uint64{in.MatchCatId, in.CatId},
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	if len(cats) < 2 {
		return error_envelope.ErrNotFound
	}

	mapCat := make(map[uint64]entity.Cat, 0)
	for _, cat := range cats {
		mapCat[cat.ID] = cat
		if cat.HasMatched {
			return error_envelope.ErrAlreadyMatch
		}
	}

	if mapCat[in.CatId].UserID == mapCat[in.MatchCatId].UserID {
		return error_envelope.ErrSameOwner
	}

	if mapCat[in.CatId].Sex.String() == mapCat[in.MatchCatId].Sex.String() {
		return error_envelope.ErrValidation
	}

	if mapCat[in.CatId].UserID != in.UserId {
		return error_envelope.ErrNotOwned
	}

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.repo.CreateMatch(ctx, tx, &entity.Match{
		UserID:      in.UserId,
		CatID:       in.CatId,
		MatchUserID: mapCat[in.MatchCatId].UserID,
		MatchCatID:  in.MatchCatId,
		Message:     in.Message,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tx.Commit()

	return nil

}

func (u *Usecase) FindMatch(ctx context.Context, in *model.FilterFindMatch) ([]model.FindMatchResponse, error) {
	logCtx := fmt.Sprintf("%T.FindMatch", u)
	var err error

	matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
		UserId:          in.UserId,
		PendingApproval: true,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	var catIds []uint64
	var userIds []uint64

	eg, _ := errgroup.WithContext(ctx)

	for _, match := range matches {
		catIds = append(catIds, match.CatID)
		catIds = append(catIds, match.MatchCatID)
		userIds = append(userIds, match.UserID)
	}

	var users []entity.User
	mapUser := make(map[uint64]model.IssuerDetail, 0)
	eg.Go(func() error {
		users, err = u.repo.FindUsers(ctx, &model.FilterFindUser{
			ID: userIds,
		})
		if err != nil {
			return err
		}

		for _, user := range users {
			mapUser[user.ID] = model.IssuerDetail{
				Email:     user.Email,
				Name:      user.Name,
				CreatedAt: user.CreatedAt.Format(constant.DefaultDateFormat),
			}
		}

		return nil
	})

	var cats []entity.Cat
	mapCat := make(map[uint64]model.FindCatResponse, 0)
	eg.Go(func() error {
		cats, err = u.repo.FindCat(ctx, &model.FilterFindCat{ID: catIds})
		if err != nil {
			return err
		}

		for _, cat := range cats {
			mapCat[cat.ID] = model.FindCatResponse{
				ID:          fmt.Sprintf("%v", cat.ID),
				Name:        cat.Name,
				Sex:         cat.Sex.String(),
				Race:        cat.Race.String(),
				ImageUrls:   cat.ImageUrls,
				AgeInMonth:  cat.Age,
				Description: cat.Description,
				HasMatched:  cat.HasMatched,
				CreatedAt:   cat.CreatedAt.Format(constant.DefaultDateFormat),
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	var resp []model.FindMatchResponse
	for _, match := range matches {
		resp = append(resp, model.FindMatchResponse{
			ID:             fmt.Sprintf("%v", match.ID),
			Message:        match.Message,
			CreatedAt:      match.CreatedAt.Format(constant.DefaultDateFormat),
			UserCatDetail:  mapCat[match.CatID],
			MatchCatDetail: mapCat[match.MatchCatID],
			IssuerDetail:   mapUser[match.UserID],
		})

	}

	return resp, nil
}

func (u *Usecase) DeleteMatch(ctx context.Context, matchId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteMatch", u)
	var err error

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
		UserId: &userId,
		ID:     []uint64{matchId},
	})

	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return error_envelope.ErrNotFound
		}
		return err
	}

	match := matches[0]

	// only issuer can delete match
	if match.UserID != userId {
		return error_envelope.ErrDeleteForbidden
	}

	if match.IsApproved || match.IsRejected {
		return error_envelope.ErrAlreadyProcessed
	}

	err = u.repo.DeleteMatch(ctx, tx, []uint64{matchId})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tx.Commit()

	return nil
}

func (u *Usecase) RejectMatch(ctx context.Context, matchId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.RejectMatch", u)
	var err error

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
		UserId: &userId,
		ID:     []uint64{matchId},
	})

	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return error_envelope.ErrNotFound
		}
		return err
	}

	match := matches[0]

	// user issued can delete match
	if match.MatchUserID != userId {
		return error_envelope.ErrDeleteForbidden
	}

	if match.IsApproved || match.IsRejected {
		return error_envelope.ErrAlreadyProcessed
	}

	err = u.repo.UpdateMatch(ctx, tx, &model.InputUpdateMatch{
		MatchId:  matchId,
		Approval: false,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tx.Commit()

	return nil
}

func (u *Usecase) ApproveMatch(ctx context.Context, matchId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.ApproveMatch", u)
	var err error

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
		UserId: &userId,
		ID:     []uint64{matchId},
	})

	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return error_envelope.ErrNotFound
		}
		return err
	}

	match := matches[0]

	// user issued can delete match
	if match.MatchUserID != userId {
		return error_envelope.ErrDeleteForbidden
	}

	if match.IsApproved || match.IsRejected {
		return error_envelope.ErrAlreadyProcessed
	}

	err = u.repo.UpdateMatch(ctx, tx, &model.InputUpdateMatch{
		MatchId:  matchId,
		Approval: true,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	go u.RemoveMatched(context.Background(), []uint64{match.MatchCatID, match.CatID})

	tx.Commit()

	return nil
}

func (u *Usecase) RemoveMatched(ctx context.Context, catIds []uint64) {
	logCtx := fmt.Sprintf("%T.ApproveMatch", u)
	var err error

	matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
		CatId: catIds,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return
	}

	var matchIds []uint64
	for _, match := range matches {
		matchIds = append(matchIds, match.ID)
	}

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.repo.DeleteMatch(ctx, tx, matchIds)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return
	}

	tx.Commit()
}
