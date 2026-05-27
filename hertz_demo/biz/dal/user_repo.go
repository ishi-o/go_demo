package dal

import (
	"context"

	"github.com/ishi-o/go_demo/hertz_demo/biz/model/entity"
	"github.com/ishi-o/go_demo/pkg/dal"

	"gorm.io/gorm"
)

type UserRepository interface {
	dal.Repository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ListByStatus(ctx context.Context, status int32, page, pageSize int) ([]*entity.User, int64, error)
	Search(ctx context.Context, keyword string, page, pageSize int) ([]*entity.User, int64, error)
}

type UserRepo struct {
	*dal.BaseRepo[entity.User]
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{
		BaseRepo: dal.NewBaseRepo[entity.User](db),
		db:       db,
	}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepo) ListByStatus(ctx context.Context, status int32, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.User{}).Where("status = ?", status)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) Search(ctx context.Context, keyword string, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.User{})
	if keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
