

package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	api "github.com/ishi-o/go_demo/hertz_demo/biz/model/api"
	"github.com/ishi-o/go_demo/hertz_demo/biz/container"
	"github.com/ishi-o/go_demo/hertz_demo/biz/service"
)


type UserHandler struct {
	userService *service.UserService
}


func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func CreateUser(ctx context.Context, c *app.RequestContext) {
	h := NewUserHandler(container.GetContainer().UserService)
	h.CreateUser(ctx, c)
}

func GetUser(ctx context.Context, c *app.RequestContext) {
	h := NewUserHandler(container.GetContainer().UserService)
	h.GetUser(ctx, c)
}

func UpdateUser(ctx context.Context, c *app.RequestContext) {
	h := NewUserHandler(container.GetContainer().UserService)
	h.UpdateUser(ctx, c)
}

func DeleteUser(ctx context.Context, c *app.RequestContext) {
	h := NewUserHandler(container.GetContainer().UserService)
	h.DeleteUser(ctx, c)
}

func ListUsers(ctx context.Context, c *app.RequestContext) {
	h := NewUserHandler(container.GetContainer().UserService)
	h.ListUsers(ctx, c)
}

func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
	var req api.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	
	user, err := h.userService.Create(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	c.JSON(consts.StatusOK, &api.CreateUserResponse{
		Id:      user.ID,
		Message: "user created successfully",
	})
}


func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	var req api.GetUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	
	user, err := h.userService.GetByID(ctx, req.Id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	c.JSON(consts.StatusOK, &api.GetUserResponse{
		User: user.ToDTO(),
	})
}


func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	var req api.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	
	user, err := h.userService.Update(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	c.JSON(consts.StatusOK, &api.UpdateUserResponse{
		User:    user.ToDTO(),
		Message: "user updated successfully",
	})
}


func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	var req api.DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	
	err := h.userService.Delete(ctx, req.Id)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	c.JSON(consts.StatusOK, &api.DeleteUserResponse{
		Success: true,
		Message: "user deleted successfully",
	})
}


func (h *UserHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	var req api.ListUsersRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	
	users, total, err := h.userService.List(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	
	userDTOs := make([]*api.User, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToDTO()
	}

	
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	
	c.JSON(consts.StatusOK, &api.ListUsersResponse{
		Users:    userDTOs,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
	})
}
