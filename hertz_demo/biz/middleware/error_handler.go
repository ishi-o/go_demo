package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ishi-o/go_demo/hertz_demo/biz/common/response"
)

func ErrorHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		if err := c.Errors.Last(); err != nil {
			response.Error(c, 500, "")
			return
		}
	}
}
