package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinsoft/it-ku/api/param/user"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"net/http"
)

var (
	UserService pb.UserService
)

func Create(c *gin.Context) {
	var regParam user.Register
	if err := c.ShouldBindJSON(&regParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	resp, err := UserService.Get(c, &pb.User{
		Email: regParam.Email,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	if resp.User.Id != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -200,
			"msg":  "邮箱已被注册",
		})
		return
	}
	if resp, err := UserService.Create(c, &pb.User{
		Name:     regParam.Name,
		Email:    regParam.Email,
		Password: regParam.Password,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"data": resp,
			"msg":  "注册成功",
		})
		return
	}
}
