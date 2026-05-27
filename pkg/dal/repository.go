package dal

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	GetByCondition(ctx context.Context, condition map[string]interface{}) (*T, error)
	List(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]T, int64, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context, condition map[string]interface{}) (int64, error)
	Exists(ctx context.Context, condition map[string]interface{}) (bool, error)
}

type BaseRepo[T any] struct {
	db *gorm.DB
}

func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{db: db}
}

func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepo[T]) GetByID(ctx context.Context, id int64) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) GetByCondition(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) List(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	db := r.db.WithContext(ctx).Model(new(T))
	for k, v := range condition {
		db = db.Where(k, v)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&entities).Error

	return entities, total, err
}

func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepo[T]) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}

func (r *BaseRepo[T]) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	var total int64
	db := r.db.WithContext(ctx).Model(new(T))
	for k, v := range condition {
		db = db.Where(k, v)
	}
	return total, db.Count(&total).Error
}

func (r *BaseRepo[T]) Exists(ctx context.Context, condition map[string]interface{}) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	for k, v := range condition {
		db = db.Where(k, v)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
