package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinsoft/it-ku/api/docs"
	"github.com/jinsoft/it-ku/api/handler/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/register", user.Create)
		apiv1.POST("/login", user.Login)
		apiv1.POST("/update", user.Update)
	}

	return r
}
