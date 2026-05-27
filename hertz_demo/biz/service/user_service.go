package service

import (
	"context"
	"errors"

	"github.com/ishi-o/go_demo/hertz_demo/biz/model/entity"
	"github.com/ishi-o/go_demo/pkg/dal"

	api "github.com/ishi-o/go_demo/hertz_demo/biz/model/api"
)

type UserService struct {
	userRepo dal.Repository[entity.User]
}

func NewUserService(userRepo dal.Repository[entity.User]) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *api.CreateUserRequest) (*entity.User, error) {
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	existingUser, err := s.userRepo.GetByCondition(ctx, map[string]interface{}{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	existingUser, err = s.userRepo.GetByCondition(ctx, map[string]interface{}{"username": req.Username})
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	user := entity.UserFromCreateRequest(req)

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *api.UpdateUserRequest) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByCondition(ctx, map[string]interface{}{"email": req.Email})
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
	}

	if req.Username != "" && req.Username != user.Username {
		existingUser, err := s.userRepo.GetByCondition(ctx, map[string]interface{}{"username": req.Username})
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("username already exists")
		}
	}

	user.FromUpdateRequest(req)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

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
