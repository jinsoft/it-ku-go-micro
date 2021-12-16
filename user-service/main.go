package main

import (
	"fmt"
	"github.com/jinsoft/it-ku/common/tracer"
	database "github.com/jinsoft/it-ku/user-service/db"
	"github.com/jinsoft/it-ku/user-service/handler"
	"github.com/jinsoft/it-ku/user-service/model"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	repository "github.com/jinsoft/it-ku/user-service/repo"
	"github.com/jinsoft/it-ku/user-service/service"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const (
	ServerName = "ik.service.user"
	EtcdAddr   = "127.0.0.1:2379"
	JaegerAddr = "127.0.0.1:6831"
)

// 启动http服务监听客户端数据采集
func prometheusBoot() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9092", nil)
		if err != nil {
			log.Fatal("listenAndServer err:", err)
		}
	}()
}

func main() {

	jaegerTracer, closer, err := tracer.NewJaegerTracer(ServerName, JaegerAddr)

	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

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
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(EtcdAddr))),
		micro.Version("latest"),
		micro.WrapHandler(
			prometheus.NewHandlerWrapper(),
			opentracing.NewHandlerWrapper(jaegerTracer),
		),
	)

	srv.Init()

	// 采集监控数据
	prometheusBoot()

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
