package main

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"log"
	database "github.com/jinsoft/it-ku/user-service/db"
	"github.com/jinsoft/it-ku/user-service/handler"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	repository "github.com/jinsoft/it-ku/user-service/repo"
)

const ServerName = "ik.user.service"

func main() {
	db, err := database.CreateConnection()
	if err != nil {
		log.Fatalf("not connect to DB: %v", err)
	}
	repo := &repository.UserRepository{Db: db}

	srv := micro.NewService(
		micro.Name(ServerName),
		micro.Version("latest"),
	)

	srv.Init()

	_ = pb.RegisterUserServiceHandler(srv.Server(), &handler.UserService{Repo: repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
