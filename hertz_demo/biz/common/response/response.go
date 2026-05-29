package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Code    int `json:"code"`
	Message any `json:"message"`
	Body    any `json:"body"`
}

func Success(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, Response{
		Code:    0,
		Message: nil,
		Body:    data,
	})
}

func SuccessEmpty(c *app.RequestContext) {
	c.JSON(consts.StatusOK, Response{
		Code:    0,
		Message: nil,
		Body:    nil,
	})
}

func Error(c *app.RequestContext, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Body:    nil,
	})
}
