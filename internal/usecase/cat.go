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
)

func (u *Usecase) CreateCat(ctx context.Context, in *model.CreateCatRequest, userId uint64) (*model.CreateCatResponse, error) {
	logCtx := fmt.Sprintf("%T.CreateCat", u)
	var err error

	inputRegister, err := validateRegisterCat(in)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	catId, err := u.repo.CreateCat(ctx, tx, inputRegister)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	tx.Commit()

	cats, err := u.repo.FindCat(ctx, &model.FilterFindCat{ID: &catId})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &model.CreateCatResponse{
		CreatedAt: cats[0].CreatedAt.Format(constant.DefaultDateFormat),
		ID:        fmt.Sprintf("%v", cats[0].ID),
	}, nil
}

func validateRegisterCat(in *model.CreateCatRequest) (*entity.Cat, error) {
	var err error

	if len(in.Name) < 1 || len(in.Name) > 30 {
		return nil, error_envelope.ErrValidation
	}

	if len(in.Sex) == 0 {
		return nil, error_envelope.ErrValidation
	}

	tSex, err := entity.StringToSex(in.Sex)
	if err != nil {
		return nil, error_envelope.ErrValidation
	}

	if len(in.Race) == 0 {
		return nil, error_envelope.ErrValidation
	}

	tRace, err := entity.StringToRace(in.Race)
	if err != nil {
		return nil, error_envelope.ErrValidation
	}

	if in.AgeInMonth <= 0 {
		return nil, error_envelope.ErrValidation
	}

	if len(in.Description) <= 0 {
		return nil, error_envelope.ErrValidation
	}

	if len(in.ImageUrls) <= 0 {
		return nil, error_envelope.ErrValidation
	}

	return &entity.Cat{
		UserID:      in.UserID,
		Name:        in.Name,
		Sex:         tSex,
		Race:        tRace,
		ImageUrls:   in.ImageUrls,
		Age:         in.AgeInMonth,
		Description: in.Description,
	}, nil
}

func (u *Usecase) DeleteCat(ctx context.Context, catId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteCat", u)
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

	err = u.repo.DeleteCat(ctx, tx, catId, userId)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tx.Commit()

	return nil
}

func (u *Usecase) UpdateCat(ctx context.Context, in *model.UpdateCatRequest) error {
	logCtx := fmt.Sprintf("%T.UpdateCat", u)
	var err error

	updateInput, err := validateUpdateCat(in)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
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

	if updateInput.Sex != nil {
		matches, err := u.repo.FindMatch(ctx, &model.FilterFindMatch{
			CatId: &in.ID,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error(ctx, logCtx, err)
			return err
		}

		if len(matches) > 0 {
			return error_envelope.ErrEmailExists
		}
	}

	err = u.repo.UpdateCat(ctx, tx, updateInput)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tx.Commit()

	return nil
}

func validateUpdateCat(in *model.UpdateCatRequest) (*model.InputUpdateCat, error) {
	var updateCat model.InputUpdateCat

	updateCat.ID = in.ID
	updateCat.UserID = in.UserID

	if in.Name != nil {
		updateCat.Name = in.Name
	}

	if in.Sex != nil {
		t, err := entity.StringToSex(*in.Sex)
		if err != nil {
			return &updateCat, error_envelope.ErrValidation
		}

		*updateCat.Sex = t.String()
	}

	if in.Race != nil {
		t, err := entity.StringToRace(*in.Race)
		if err != nil {
			return &updateCat, error_envelope.ErrValidation
		}

		*updateCat.Race = t.String()
	}

	if in.ImageUrls != nil || len(in.ImageUrls) > 0 {
		updateCat.ImageUrls = in.ImageUrls
	}

	if in.AgeInMonth != nil {
		updateCat.Age = in.AgeInMonth
	}

	if in.Description != nil {
		updateCat.Description = in.Description
	}

	return nil, nil
}

func (u *Usecase) FindCat(ctx context.Context, in *model.FilterFindCat) ([]model.FindCatResponse, error) {
	logCtx := fmt.Sprintf("%T.DeleteCat", u)
	var err error

	cats, err := u.repo.FindCat(ctx, in)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	var resp []model.FindCatResponse
	for _, cat := range cats {
		resp = append(resp, model.FindCatResponse{
			ID:          fmt.Sprintf("%v", cat.ID),
			Name:        cat.Name,
			Sex:         cat.Sex.String(),
			Race:        cat.Race.String(),
			ImageUrls:   cat.ImageUrls,
			AgeInMonth:  cat.Age,
			Description: cat.Description,
			HasMatched:  cat.HasMatched,
			CreatedAt:   cat.CreatedAt.Format(constant.DefaultDateFormat),
		})
	}

	return resp, nil
}
