package domain

import (
	"errors"
	"net/http"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")

	ErrEmailAlreadyUsed = errors.New("Email already used")

	ErrForbidden = errors.New("Forbidden")

	ErrLoginFailed = errors.New("Invalid Email or Password")
)

// unfinished
func GetStatusCode(err error) int {
	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrBadParamInput:
		return http.StatusBadRequest
	case ErrEmailAlreadyUsed:
		return http.StatusBadRequest
	case ErrForbidden:
		return http.StatusForbidden
	case ErrLoginFailed:
		return http.StatusUnauthorized
	default:
		return http.StatusOK
	}
}
