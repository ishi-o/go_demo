package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/ishi-o/go_demo/hertz_demo/biz/common/response"
	"github.com/ishi-o/go_demo/pkg/log"
	"go.uber.org/zap"
)

func ErrorHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		if err := c.Errors.Last(); err != nil {
			log.Warn(
				"error happened",
				zap.String("error", fmt.Sprintf("%+v", err.Err)),
			)
			response.Error(c, consts.StatusInternalServerError, "")
			return
		}
	}
}
