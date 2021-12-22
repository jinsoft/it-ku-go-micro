package main

import (
	r "github.com/jinsoft/it-ku/api/router"
	serviceclient "github.com/jinsoft/it-ku/api/service-client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"log"
	"os"
)

var (
	ServerName = "ik.web.api"
	//EtcdAddr   = "127.0.0.1:2379"
	EtcdAddr = os.Getenv("MICRO_REGISTRY_ADDRESS")
)

func init() {
	serviceclient.RegisterService()
}

// @title IK API
// @version 1.0
// @description This is a sample api

// @contact.name hhh
// @contact.url https://www.ainiok.com
// @contact.email job@ainiok.com

// @host localhost:8088
// @BasePath /v1
func main() {

	g := r.NewRouter()

	srv := web.NewService(
		web.Name(ServerName),
		web.Address(":8088"),
		web.Registry(etcd.NewRegistry(
			registry.Addrs(EtcdAddr))),
		web.Handler(g),
	)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
