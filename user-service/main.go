package main

import (
	"fmt"
	database "github.com/jinsoft/it-ku/user-service/db"
	"github.com/jinsoft/it-ku/user-service/handler"
	"github.com/jinsoft/it-ku/user-service/model"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	repository "github.com/jinsoft/it-ku/user-service/repo"
	"github.com/jinsoft/it-ku/user-service/service"
	"github.com/micro/go-micro/v2"
	"log"
)

const ServerName = "ik.service.user"

func main() {
	db, err := database.CreateConnection()
	if err != nil {
		log.Fatalf("not connect to DB: %v", err)
	}

	// 每次启动服务时都会检查，如果数据表不存在则创建，已存在检查是否有修改
	db.AutoMigrate(&model.User{})

	repo := &repository.UserRepository{db}
	//token := &service.TokenService{repo}

	srv := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	srv.Init()

	userHandler := &handler.UserService{
		Repo:      repo,
		Token:     &service.TokenService{repo},
		ResetRepo: &repository.PasswordResetRepository{db},
		PubSub:    srv.Server().Options().Broker,
	}

	_ = pb.RegisterUserServiceHandler(srv.Server(), userHandler)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
