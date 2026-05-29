package errors

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	pkgerr "github.com/pkg/errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInternalServer    = errors.New("internal server error")
)

func WrapWithStack(c *app.RequestContext, err error) error {
	return c.Error(pkgerr.WithStack(err))
}
