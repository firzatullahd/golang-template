package service

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/entity/converter"
	"github.com/firzatullahd/golang-template/internal/user/model"
	customerror "github.com/firzatullahd/golang-template/internal/user/model/error"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, in model.RegisterRequest) (*model.RegisterResponse, error) {
	// logCtx := fmt.Sprintf("%T.Register", u)

	if err := validateRegister(&in); err != nil {
		return nil, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		// logger.
		return nil, err
	}
	tx, err := s.repo.WithTransaction()
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// checkUser, err := s.repo.FindUsers(ctx, &model.FilterFindUser{Email: &in.Username})
	// if err != nil && !errors.Is(err, sql.ErrNoRows) {
	// 	return nil, err
	// }

	// if checkUser != nil {
	// 	return nil, customerror.ErrEmailExists
	// }

	userId, err := s.repo.CreateUser(ctx, tx, converter.RegisterRequestToEntity(in, password))
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, err
	}

	// TODO: handle error duplicate unique key

	accessToken, err := s.generateAccessToken(userId, in.Username)
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, err
	}

	tx.Commit()

	return &model.RegisterResponse{
		Username:    in.Username,
		Name:        in.Name,
		AccessToken: accessToken,
	}, nil
}

func validateRegister(in *model.RegisterRequest) error {
	_, err := mail.ParseAddress(in.Username)
	if err != nil {
		return customerror.ErrValidation
	}

	if len(in.Name) < 5 || len(in.Name) > 50 {
		return customerror.ErrValidation
	}

	if len(in.Password) < 5 || len(in.Password) > 15 {
		return customerror.ErrValidation
	}

	return nil
}

func (s *Service) Login(ctx context.Context, in model.AuthRequest) (*model.AuthResponse, error) {
	// logCtx := fmt.Sprintf("%T.Login", u)
	user, err := s.repo.FindUser(ctx, &model.FilterFindUser{
		Username: &in.Username,
	})
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, customerror.ErrWrongPassword
	}

	accessToken, err := s.generateAccessToken(user.ID, in.Username)
	if err != nil {
		// logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &model.AuthResponse{
		Name:        user.Name,
		Username:    user.Username,
		AccessToken: accessToken,
	}, nil

}

func (s *Service) generateAccessToken(userId uint64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaim{
		UserData: model.UserData{
			UserID:   userId,
			Username: email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	})

	return token.SignedString([]byte(s.conf.JWTSecretKey))
}
