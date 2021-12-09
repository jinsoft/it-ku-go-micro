package service_client

import (
	"github.com/jinsoft/it-ku/api/handler/user"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"github.com/micro/go-micro/v2"
)

func RegisterService() {
	app := micro.NewService(
		micro.Name("ik.client.api"),
		micro.Version("latest"),
	)

	cli := app.Client()
	user.Srv = pb.NewUserService("ik.service.user", cli)
}
