package main

import (
	r "github.com/jinsoft/it-ku/api/router"
	serviceclient "github.com/jinsoft/it-ku/api/service-client"
	"github.com/micro/go-micro/v2/web"
	"log"
)

const ServerName = "ik.web.api"

func init() {
	serviceclient.RegisterService()
}
func main() {

	g := r.NewRouter()

	srv := web.NewService(
		web.Name(ServerName),
		web.Address(":8088"),
		web.Handler(g),
	)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
