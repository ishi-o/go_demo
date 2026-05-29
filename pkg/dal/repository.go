package dal

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	GetByCondition(ctx context.Context, condition map[string]any) (*T, error)
	List(ctx context.Context, condition map[string]any, page, pageSize int) ([]T, int64, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, condition map[string]any) (int64, error)
	Exists(ctx context.Context, condition map[string]any) (bool, error)
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

func (r *BaseRepo[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) GetByCondition(ctx context.Context, condition map[string]any) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) List(ctx context.Context, condition map[string]any, page, pageSize int) (entities []T, total int64, err error) {
	db := r.db.WithContext(ctx).Model(new(T))
	for k, v := range condition {
		db = db.Where(k, v)
	}

	offset := (page - 1) * pageSize
	err = db.Offset(offset).Limit(pageSize).Find(&entities).Error
	total = int64(len(entities))
	return
}

func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepo[T]) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(new(T)).Error
}

func (r *BaseRepo[T]) Count(ctx context.Context, condition map[string]any) (int64, error) {
	var total int64
	db := r.db.WithContext(ctx).Model(new(T))
	for k, v := range condition {
		db = db.Where(k, v)
	}
	return total, db.Count(&total).Error
}

func (r *BaseRepo[T]) Exists(ctx context.Context, condition map[string]any) (bool, error) {
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
