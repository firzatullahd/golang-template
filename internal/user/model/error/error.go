package customerror

import (
	"fmt"
	"net/http"
)

var (
	ErrEmailExists   = fmt.Errorf("email already exists")
	ErrValidation    = fmt.Errorf("please use valid email")
	ErrNotFound      = fmt.Errorf("not found")
	ErrWrongPassword = fmt.Errorf("wrong password")
	ErrUnknown       = fmt.Errorf("error unknown")
	ErrUnauthorized  = fmt.Errorf("unauthorized")
)

var mapErrorCode = map[error]int{
	ErrEmailExists:   http.StatusConflict,
	ErrValidation:    http.StatusBadRequest,
	ErrNotFound:      http.StatusNotFound,
	ErrWrongPassword: http.StatusBadRequest,
	ErrUnknown:       http.StatusInternalServerError,
	ErrUnauthorized:  http.StatusUnauthorized,
}

func ParseError(err error) (code int, err2 error) {
	if err != nil {
		if errCode, ok := mapErrorCode[err]; ok {
			code = errCode
			return code, err2
		} else {
			return http.StatusInternalServerError, ErrUnknown
		}
	}
	return http.StatusOK, nil
}
