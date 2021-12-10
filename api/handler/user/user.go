package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinsoft/it-ku/api/param/user"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"log"
	"net/http"
	"strconv"
)

var (
	Srv pb.UserService
)

func Update(c *gin.Context) {
	var userParam user.User
	if err := c.ShouldBindJSON(&userParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	uid := strconv.FormatInt(userParam.Id, 10)
	resp, err := Srv.Update(context.TODO(), &pb.User{
		Id:       uid,
		Name:     userParam.Name,
		Email:    userParam.Email,
		Password: userParam.Password,
		Status:   userParam.Status,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": resp,
		"msg":  "更新成功",
	})
	return
}

func Login(c *gin.Context) {
	// 先实现邮箱登录， 后续再扩展成手机号和验证码登录
	var login user.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	token, err := Srv.Auth(context.TODO(), &pb.User{
		Email:    login.Email,
		Password: login.Password,
	})
	if err != nil {
		log.Printf("用户登录失败:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -200,
			"msg":  "登录失败",
		})
		return
	}
	log.Printf("用户登录成功: %s", token.Token)
	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"msg":   "登录成功",
		"token": token.Token,
	})
	return
}

func Create(c *gin.Context) {
	var regParam user.Register
	if err := c.ShouldBindJSON(&regParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	resp, err := Srv.Get(c, &pb.User{
		Email: regParam.Email,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	if resp.User != nil && resp.User.Id != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -200,
			"msg":  "邮箱已被注册",
		})
		return
	}

	if resp, err := Srv.Create(c, &pb.User{
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
