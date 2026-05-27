package service

import (
	"context"

	"github.com/ishi-o/go_demo/hertz_demo/biz/common/errors"
	"github.com/ishi-o/go_demo/hertz_demo/biz/model/entity"
	"github.com/ishi-o/go_demo/pkg/dal"

	api "github.com/ishi-o/go_demo/hertz_demo/biz/model/api"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo dal.Repository[entity.User]
}

func NewUserService(userRepo dal.Repository[entity.User]) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *api.CreateUserRequest) error {
	user := entity.UserFromCreateRequest(req)
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *api.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return err
	}

	user.FromUpdateRequest(req)

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context, req *api.ListUsersRequest) ([]*entity.User, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	condition := map[string]interface{}{}
	if req.Status > 0 {
		condition["status"] = int32(req.Status)
	}

	users, total, err := s.userRepo.List(ctx, condition, int(page), int(pageSize))
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	return result, total, nil
}
