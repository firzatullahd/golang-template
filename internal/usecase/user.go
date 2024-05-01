package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/firzatullahd/cats-social-api/internal/entity"
	"github.com/firzatullahd/cats-social-api/internal/model"
	"github.com/firzatullahd/cats-social-api/internal/utils/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) Register(ctx context.Context, in *model.RegisterRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)
	var err error
	// validate payload

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

	accessToken, err := u.generateAccessToken(fmt.Sprintf("%v", userId))
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

func (u *Usecase) Login(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)
	user, err := u.repo.FindUser(ctx, in.Email)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email/password salah")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, errors.New("email/password salah")
	}

	accessToken, err := u.generateAccessToken(fmt.Sprintf("%v", user.ID))
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

func (u *Usecase) generateAccessToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaim{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	})

	return token.SignedString(model.JWT_SIGNATURE_KEY)
}
