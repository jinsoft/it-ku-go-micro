package service_client

import (
	"github.com/jinsoft/it-ku/api/handler/user"
	"github.com/jinsoft/it-ku/common/tracer"
	"github.com/jinsoft/it-ku/common/wrapper/breaker/hystrix"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"log"
)

const (
	ServerName = "ik.client.api"
	JaegerAddr = "192.168.1.32:6831"
)

func RegisterService() {

	jaegerTracer, closer, err := tracer.NewJaegerTracer(ServerName, JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	app := micro.NewService(
		micro.Name("ik.client.api"),
		micro.Version("latest"),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	cli := app.Client()
	user.Srv = pb.NewUserService("ik.service.user", cli)
}
