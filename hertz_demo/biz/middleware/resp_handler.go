package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ishi-o/go_demo/hertz_demo/biz/common/response"
	"github.com/ishi-o/go_demo/pkg/log"
	"go.uber.org/zap"
)

func RespHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		if c.Errors.Last() != nil {
			return
		}
		if data, exists := c.Get("response_data"); exists {
			log.Info("response with body", zap.Any("body", data))
			response.Success(c, data)
		} else {
			log.Info("empty response")
			response.SuccessEmpty(c)
		}
	}
}
