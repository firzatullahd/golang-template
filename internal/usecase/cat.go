package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/firzatullahd/cats-social-api/internal/utils/constant"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
)

func (u *Usecase) CreateCat(ctx context.Context, in *model.CreateCatRequest, userId uint64) (*model.CreateCatResponse, error) {
	logCtx := fmt.Sprintf("%T.CreateCat", u)
	var err error

	if err := validateCat(in); err != nil {
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

	catId, err := u.repo.CreateCat(ctx, tx, &entity.Cat{
		ID:          userId,
		UserID:      userId,
		Name:        in.Name,
		Sex:         in.Sex,
		Race:        in.Race,
		ImageUrls:   in.ImageUrls,
		Age:         in.AgeInMonth,
		Description: in.Description,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	tx.Commit()

	return &model.CreateCatResponse{
		CreatedAt: time.Now().Format(constant.DefaultDateFormat),
		ID:        catId,
	}, nil
}

func validateCat(in *model.CreateCatRequest) error {
	var err error
	if len(in.Name) < 1 || len(in.Name) > 30 {
		return error_envelope.ErrValidation
	}

	if len(in.SexStr) == 0 {
		return error_envelope.ErrValidation
	}

	in.Sex, err = entity.StringToSex(in.SexStr)
	if err != nil {
		return error_envelope.ErrValidation
	}

	if len(in.RaceStr) == 0 {
		return error_envelope.ErrValidation
	}

	in.Race, err = entity.StringToRace(in.RaceStr)
	if err != nil {
		return error_envelope.ErrValidation
	}

	if in.AgeInMonth <= 0 {
		return error_envelope.ErrValidation
	}

	if len(in.Description) <= 0 {
		return error_envelope.ErrValidation
	}

	if len(in.ImageUrls) <= 0 {
		return error_envelope.ErrValidation
	}

	return nil
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

func (u *Usecase) UpdateCat(ctx context.Context, in *model.CreateCatRequest, catId, userId uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteCat", u)
	var err error

	err = validateCat(in)
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

	err = u.repo.UpdateCat(ctx, tx, &entity.Cat{
		ID:          catId,
		UserID:      userId,
		Name:        in.Name,
		Sex:         in.Sex,
		Race:        in.Race,
		ImageUrls:   in.ImageUrls,
		Age:         in.AgeInMonth,
		Description: in.Description,
		HasMatched:  false,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	// todo if already request match, return error when update sex

	tx.Commit()

	return nil
}

func (u *Usecase) FindCat(ctx context.Context, in *model.FilterFindCat) error {
	logCtx := fmt.Sprintf("%T.DeleteCat", u)
	var err error

	_, err = u.repo.FindCat(ctx, in)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}
