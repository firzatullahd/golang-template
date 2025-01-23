package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/entity"
	"github.com/firzatullahd/golang-template/internal/user/entity/converter"
	"github.com/firzatullahd/golang-template/internal/user/model"
	customerror "github.com/firzatullahd/golang-template/internal/user/model/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, in model.RegisterRequest) (*model.RegisterResponse, error) {
	logCtx := fmt.Sprintf("%T.Register", s)

	if err := validateRegister(&in); err != nil {
		return nil, fmt.Errorf("%s %w", logCtx, err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logCtx, err)
	}
	tx, err := s.repo.WithTransaction()
	if err != nil {
		return nil, fmt.Errorf("%s %w", logCtx, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userId, err := s.repo.CreateUser(ctx, tx, converter.RegisterRequestToEntity(in, password))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if strings.EqualFold(string(pqErr.Code), "23505") {
				return nil, customerror.ErrUsernameExists
			}
		}

		return nil, fmt.Errorf("%s %w", logCtx, err)
	}

	accessToken, err := s.generateAccessToken(userId, in.Username)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logCtx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s %w", logCtx, err)
	}

	return &model.RegisterResponse{
		Username:    in.Username,
		Name:        in.Name,
		AccessToken: accessToken,
		UserID:      userId,
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
	user, err := s.repo.FindUser(ctx, &model.FilterFindUser{
		Username: &in.Username,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, customerror.ErrWrongPassword
	}

	accessToken, err := s.generateAccessToken(user.ID, in.Username)
	if err != nil {
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

	if err := s.repo.UpdateUser(ctx, tx, &model.FilterFindUser{Username: &username}, map[string]interface{}{
		"state": entity.UserStatePending,
	}); err != nil {
		return err
	}

	if err := s.emailClient.SendEmail(ctx, model.EmailPayload{
		Email:            username,
		Name:             user.Name,
		VerificationCode: verificationCode,
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
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
	key := fmt.Sprintf(model.VerificationCounterPrefix, username)

	val, err := s.redisConn.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	if val == 0 {
		fmt.Println("LALA", s.time.UntilMidnight())
		err = s.redisConn.Set(ctx, key, 1, s.time.UntilMidnight()).Err()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if val >= model.VerificationMaxAttempt {
		return false, customerror.ErrTooManyRequests
	}

	err = s.redisConn.Incr(ctx, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Verification(ctx context.Context, username, code string) error {

	val, err := s.redisConn.Get(ctx, fmt.Sprintf(model.VerificationPrefix, username)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if val != code {
		return customerror.ErrInvalidVerificationCode
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

	if err := s.repo.UpdateUser(ctx, tx, &model.FilterFindUser{Username: &username}, map[string]interface{}{
		"state": entity.UserStateVerified,
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
