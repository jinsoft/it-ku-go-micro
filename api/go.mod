module github.com/jinsoft/it-ku/api

go 1.13

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/jinsoft/it-ku/user-service v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro/v2 v2.9.1
)

replace github.com/jinsoft/it-ku/user-service => ./../user-service
