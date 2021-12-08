package handler

import (
	"github.com/gin-gonic/gin"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
)

var (
	UserService pb.UserService
)

func Create(c *gin.Context) {
	req := new(pb.User)
	if err := c.BindQuery(req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "bad param",
		})
		return
	}
	if resp, err := UserService.Create(c, req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": "200",
			"data": resp,
		})
	}
}
