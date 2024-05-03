package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) Register(ctx context.Context, in *model.RegisterRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)

	if err := validateRegister(in); err != nil {
		return nil, err
	}

	// todo cost config
	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
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

	checkUser, err := u.repo.FindUser(ctx, in.Email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	if checkUser != nil {
		return nil, error_envelope.ErrEmailExists
	}

	user := entity.User{
		Email:    in.Email,
		Password: string(password),
		Name:     in.Name,
	}
	userId, err := u.repo.CreateUser(ctx, tx, &user)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	accessToken, err := u.generateAccessToken(userId, in.Email)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	tx.Commit()

	return &model.AuthResponse{
		Email:       in.Email,
		Name:        in.Name,
		AccessToken: accessToken,
	}, nil
}

func validateRegister(in *model.RegisterRequest) error {
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		return error_envelope.ErrValidation
	}

	if len(in.Name) < 5 || len(in.Name) > 50 {
		return error_envelope.ErrValidation
	}

	if len(in.Password) < 5 || len(in.Password) > 15 {
		return error_envelope.ErrValidation
	}

	return nil
}

func (u *Usecase) Login(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)
	user, err := u.repo.FindUser(ctx, in.Email)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_envelope.ErrNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, error_envelope.ErrWrongPass
	}

	accessToken, err := u.generateAccessToken(user.ID, in.Email)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &model.AuthResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil

}

func (u *Usecase) generateAccessToken(userId uint64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaim{
		UserData: model.UserData{
			ID:    userId,
			Email: email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	})

	return token.SignedString(model.JWT_SIGNATURE_KEY)
}
