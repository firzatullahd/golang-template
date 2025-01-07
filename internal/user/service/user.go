package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/entity"
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

	if err := tx.Commit(); err != nil {
		return nil, err
	}

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

func (s *Service) generateAccessToken(userId uint64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaim{
		UserData: model.UserData{
			UserID:   userId,
			Username: username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	})

	return token.SignedString([]byte(s.conf.JWTSecretKey))
}

func (s *Service) InitialVerification(ctx context.Context, username string) error {

	user, err := s.repo.FindUser(ctx, &model.FilterFindUser{
		Username: &username,
	})
	if err != nil {
		return err
	}

	allowed, err := s.allowInitialVerification(ctx, username)
	if err != nil {
		return err
	}

	if !allowed {
		return customerror.ErrTooManyRequests
	}

	switch user.State {
	case entity.UserStateVerified:
		return customerror.ErrAlreadyVerified
	case entity.UserStateDeleted:
		return customerror.ErrNotFound
	}

	verificationCode, err := s.generateVerificationCode()
	if err != nil {
		return err
	}

	if err := s.redisConn.Set(ctx, fmt.Sprintf(model.VerificationPrefix, username), verificationCode, model.VerificationTTL).Err(); err != nil {
		return err
	}

	tx, err := s.repo.WithTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := s.repo.UpdateUser(ctx, tx, user.ID, map[string]interface{}{
		"state": entity.UserStatePending,
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := s.emailClient.SendEmail(ctx, model.EmailPayload{
		Email:            username,
		Name:             user.Name,
		VerificationCode: verificationCode,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) generateVerificationCode() (string, error) {
	const min, max = 100000, 999999

	rangeSize := max - min + 1

	num, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		return "", err
	}

	otp := int(num.Int64()) + min

	return fmt.Sprintf("%06d", otp), nil
}

func (s *Service) allowInitialVerification(ctx context.Context, username string) (bool, error) {
	val, err := s.redisConn.Get(ctx, fmt.Sprintf(model.VerificationCounterPrefix, username)).Int()
	if err != nil {
		return false, err
	}

	if val == 0 {
		// set initial counter
		return true, s.redisConn.Set(ctx, fmt.Sprintf(model.VerificationCounterPrefix, username), 1, s.time.UntilMidnight()).Err()
	}

	if val >= model.VerificationMaxAttempt {
		return false, customerror.ErrTooManyRequests
	}

	// increment counter
	return true, s.redisConn.Incr(ctx, fmt.Sprintf(model.VerificationCounterPrefix, username)).Err()
}

func (s *Service) Verify(ctx context.Context, username, code string) error {

	val, err := s.redisConn.Get(ctx, fmt.Sprintf(model.VerificationPrefix, username)).Result()
	if err != nil {
		return err
	}

	if val != code {
		return customerror.ErrInvalidVerificationCode
	}

	return nil
}

func (s *Service) DoKyc(ctx context.Context) {

}
