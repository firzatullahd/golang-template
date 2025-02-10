package customerror

import (
	"errors"
	"net/http"
)

var (
	ErrUsernameExists          = errors.New("username already exists")
	ErrValidationUsername      = errors.New("please use valid email")
	ErrValidationName          = errors.New("name must be between 5-50 characters")
	ErrValidationPassword      = errors.New("password must be between 5-15 characters")
	ErrNotFound                = errors.New("not found")
	ErrWrongPassword           = errors.New("wrong password")
	ErrUnknown                 = errors.New("error unknown")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrAlreadyVerified         = errors.New("user already verified")
	ErrNoResourceUpdated       = errors.New("no resource updated")
	ErrTooManyRequests         = errors.New("too many requests")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
)

var mapErrorCode = map[error]int{
	ErrUsernameExists:          http.StatusConflict,
	ErrValidationUsername:      http.StatusBadRequest,
	ErrValidationName:          http.StatusBadRequest,
	ErrValidationPassword:      http.StatusBadRequest,
	ErrNotFound:                http.StatusNotFound,
	ErrWrongPassword:           http.StatusBadRequest,
	ErrUnknown:                 http.StatusInternalServerError,
	ErrUnauthorized:            http.StatusUnauthorized,
	ErrAlreadyVerified:         http.StatusConflict,
	ErrTooManyRequests:         http.StatusTooManyRequests,
	ErrInvalidVerificationCode: http.StatusBadRequest,
}

func ParseError(err error) (int, error) {
	if err != nil {
		if code, ok := mapErrorCode[err]; ok {
			return code, err
		} else {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}
