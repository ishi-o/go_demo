package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ishi-o/go_demo/hertz_demo/biz/common/response"
)

func RespHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		if data, exists := c.Get("response_data"); exists {
			response.Success(c, data)
		} else {
			response.SuccessEmpty(c)
		}
	}
}
