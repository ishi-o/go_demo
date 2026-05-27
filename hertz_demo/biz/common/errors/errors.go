package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInternalServer    = errors.New("internal server error")
)
