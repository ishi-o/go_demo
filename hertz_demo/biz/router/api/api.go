package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	handler "github.com/ishi-o/go_demo/hertz_demo/biz/handler/api"
	"github.com/ishi-o/go_demo/hertz_demo/biz/middleware"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

func Register(r *server.Hertz) {
	v1Group := r.Group("/api/v1")
	v1Group.Use(
		middleware.Logger(),
		middleware.CORS(),
		middleware.Auth(),
		middleware.ErrorHandlerMiddleware(),
		middleware.RespHandlerMiddleware(),
	)

	userGroup := v1Group.Group("/users")
	{
		userGroup.POST("", handler.CreateUser)
		userGroup.GET("/:id", handler.GetUser)
		userGroup.PUT("/:id", handler.UpdateUser)
		userGroup.DELETE("/:id", handler.DeleteUser)
		userGroup.GET("", handler.ListUsers)
	}
}
