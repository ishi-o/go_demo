package main

import (
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/ishi-o/go_demo/hertz_demo/biz/container"
	"github.com/ishi-o/go_demo/hertz_demo/biz/model/entity"
	"google.golang.org/protobuf/proto"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/ishi-o/go_demo/pkg/config"
	"github.com/ishi-o/go_demo/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	err = log.Init(&cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer log.Sync()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect database: %v", err))
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to migrate database: %v", err))
	}
	container.Init(db)

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.Server.Port)),
		server.WithCustomValidatorFunc(func(_ *protocol.Request, req any) error {
			if msg, ok := req.(proto.Message); ok {
				return protovalidate.Validate(msg)
			}
			return nil
		}),
	)
	register(h)

	h.Spin()
}
