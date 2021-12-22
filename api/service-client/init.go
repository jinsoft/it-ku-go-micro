package service_client

import (
	"github.com/jinsoft/it-ku/api/handler/user"
	"github.com/jinsoft/it-ku/common/tracer"
	"github.com/jinsoft/it-ku/common/wrapper/breaker/hystrix"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"log"
	"os"
)

var (
	ServerName = "ik.client.api"
	JaegerAddr = os.Getenv("MICRO_TRACE_SERVER")
	EtcdAddr   = os.Getenv("MICRO_REGISTRY_ADDRESS")
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
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(EtcdAddr))),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	cli := app.Client()
	user.Srv = pb.NewUserService("ik.service.user", cli)
}
