module github.com/jinsoft/it-ku/user-service

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/jinsoft/it-ku/common v0.0.0-20211215105525-346396e22c87 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.9.1 // indirect
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1 // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/mysql v1.2.1
	gorm.io/gorm v1.22.4
)
