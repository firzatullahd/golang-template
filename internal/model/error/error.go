package error_envelope

import (
	"fmt"
	"net/http"
)

var (
	ErrEmailExists  = fmt.Errorf("email already exists")
	ErrValidation   = fmt.Errorf("validation error")
	ErrNotFound     = fmt.Errorf("not found")
	ErrWrongPass    = fmt.Errorf("wrong password")
	ErrUnknown      = fmt.Errorf("error unknown")
	ErrUnauthorized = fmt.Errorf("unauthorized")
)

var mapErrorCode = map[error]int{
	ErrEmailExists:  http.StatusConflict,
	ErrValidation:   http.StatusBadRequest,
	ErrNotFound:     http.StatusNotFound,
	ErrWrongPass:    http.StatusBadRequest,
	ErrUnknown:      http.StatusInternalServerError,
	ErrUnauthorized: http.StatusUnauthorized,
}

func ParseError(err error) (code int, errMsg string) {
	if err != nil {
		if errCode, ok := mapErrorCode[err]; ok {
			code = errCode
			return code, err.Error()
		} else {
			return http.StatusInternalServerError, ErrUnknown.Error()
		}
	}
	return http.StatusOK, ""
}
