package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ishi-o/go_demo/hertz_demo/biz/common/errors"
	"github.com/ishi-o/go_demo/hertz_demo/biz/container"
	api "github.com/ishi-o/go_demo/hertz_demo/biz/model/api"
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
	if err := c.BindAndValidate(&req); err != nil {
		errors.WrapWithStack(c, errors.ErrInvalidRequest)
		return
	}

	if err := h.userService.Create(ctx, &req); err != nil {
		errors.WrapWithStack(c, err)
		return
	}

	c.Set("response_data", nil)
}

func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	var req api.GetUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.WrapWithStack(c, errors.ErrInvalidRequest)
		return
	}

	user, err := h.userService.GetByID(ctx, req.Id)
	if err != nil {
		errors.WrapWithStack(c, err)
		return
	}

	c.Set("response_data", user.ToDTO())
}

func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	var req api.UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.WrapWithStack(c, errors.ErrInvalidRequest)
		return
	}

	if err := h.userService.Update(ctx, &req); err != nil {
		errors.WrapWithStack(c, err)
		return
	}

	c.Set("response_data", nil)
}

func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	var req api.DeleteUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.WrapWithStack(c, errors.ErrInvalidRequest)
		return
	}

	if err := h.userService.Delete(ctx, req.Id); err != nil {
		errors.WrapWithStack(c, err)
		return
	}

	c.Set("response_data", nil)
}

func (h *UserHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	var req api.ListUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.WrapWithStack(c, errors.ErrInvalidRequest)
		return
	}

	users, total, err := h.userService.List(ctx, &req)
	if err != nil {
		errors.WrapWithStack(c, err)
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

	c.Set("response_data", &api.ListUsersResponse{
		Users:    userDTOs,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
	})
}
