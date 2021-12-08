package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinsoft/it-ku/api/handler"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/register", handler.Create)
	}

	return r
}
