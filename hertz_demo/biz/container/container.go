package container

import (
	"github.com/ishi-o/go_demo/hertz_demo/biz/model/entity"
	"github.com/ishi-o/go_demo/hertz_demo/biz/service"
	"github.com/ishi-o/go_demo/pkg/dal"

	"gorm.io/gorm"
)

type Container struct {
	DB          *gorm.DB
	UserRepo    dal.Repository[entity.User]
	UserService *service.UserService
}

var globalContainer *Container

func Init(db *gorm.DB) {
	globalContainer = &Container{
		DB:          db,
		UserRepo:    dal.NewBaseRepo[entity.User](db),
		UserService: service.NewUserService(dal.NewBaseRepo[entity.User](db)),
	}
}

func GetContainer() *Container {
	return globalContainer
}
