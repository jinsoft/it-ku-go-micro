package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinsoft/it-ku/user-service/model"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"github.com/jinsoft/it-ku/user-service/repo"
	"github.com/jinsoft/it-ku/user-service/service"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

const (
	passwordResetTopic = "password.reset"
	aciveAccountTopic  = "account.active"
)

type UserService struct {
	Repo      repo.Repository
	Token     service.Authable
	ResetRepo repo.PasswordResetInterface
	PubSub    broker.Broker
}

func (srv *UserService) GetById(ctx context.Context, user *pb.User, response *pb.Response) error {
	panic("implement me")
}

func (srv *UserService) CreatePasswordReset(ctx context.Context, req *pb.PasswordReset, res *pb.PasswordResetResponse) error {
	if req.Email == "" {
		return errors.New("", "邮箱不能为空", http.StatusBadRequest)
	}
	resetModel := new(model.PasswordReset)
	passwordReset, _ := resetModel.ToORM(req)
	if err := srv.ResetRepo.Create(req); err != nil {
		return err
	}
	if passwordReset != nil {
		req, _ = passwordReset.ToProtobuf()
		if err := srv.publishEvent(req); err != nil {
			return err
		}
		res.PasswordReset = req
	}
	return nil
}

func (srv *UserService) ValidatePasswordResetToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	if req.Token == "" {
		return errors.New("", "Token信息不能为空", http.StatusBadRequest)
	}
	_, err := srv.ResetRepo.GetByToken(req.Token)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("", "数据库查询异常", http.StatusBadRequest)
	}
	if err == gorm.ErrRecordNotFound {
		res.Vlaid = false
	} else {
		res.Vlaid = true
	}
	return nil
}

func (srv *UserService) DeletePasswordReset(ctx context.Context, reset *pb.PasswordReset, response *pb.PasswordResetResponse) error {
	panic("implement me")
}

func (srv *UserService) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	var userModel *model.User
	var err error
	if req.Id != "" {
		id, _ := strconv.ParseUint(req.Id, 10, 64)
		userModel, err = srv.Repo.Get(uint(id))
	} else if req.Email != "" {
		userModel, err = srv.Repo.GetByEmail(req.Email)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if userModel != nil {
		res.User, _ = userModel.ToProtobuf()
	}
	return nil
}

func (srv *UserService) Update(ctx context.Context, req *pb.User, res *pb.Response) error {
	if req.Id == "" {
		return errors.New("", "ID 不能为空", 402)
	}
	if req.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		req.Password = string(hashedPass)
	}
	id, _ := strconv.ParseUint(req.Id, 10, 64)
	user, err := srv.Repo.Get(uint(id))
	if err != nil && err != gorm.ErrRecordNotFound {
		// 记录不存在
		return errors.New("", "用户不存在", 402)
	}
	uptUser, _ := user.ToORM(req)
	// todo: 这里有个问题， 更新的时候会同时更新created_at 的值，  UPDATE `users` SET `created_at`='0000-00-00 00:00:00',`updated_at`='2021-12-10 16:38:58.391' 暂未解决
	if err := srv.Repo.Update(uptUser); err != nil {
		return err
	}
	return nil
}

func (srv *UserService) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.Repo.GetAll()
	if err != nil {
		return err
	}
	userItems := make([]*pb.User, len(users))
	for index, user := range users {
		userItem, _ := user.ToProtobuf()
		userItems[index] = userItem
	}
	res.Users = userItems
	return nil
}

func (srv *UserService) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	// 对密码进行哈希加密
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPass)
	userModel := &model.User{}
	user, _ := userModel.ToORM(req)
	if err := srv.Repo.Create(user); err != nil {
		return err
	}

	// 注册成功, 发送激活账号邮件
	createStuct := map[string]string{
		"email":    req.Email,
		"password": req.Password,
	}
	body, _ := json.Marshal(createStuct)
	msg := &broker.Message{
		Header: map[string]string{
			"email": req.Email,
		},
		Body: body,
	}
	fmt.Println("*****************pub")
	if err := srv.PubSub.Publish(aciveAccountTopic, msg); err != nil {
		log.Printf("[pub] create account pub failed: %v", err)
	}
	res.User = req
	return nil
}

func (srv *UserService) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.Repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}
	// 生成 jwt token
	token, err := srv.Token.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *UserService) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := srv.Token.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.ID == 0 {
		return errors.New("", "无效用户", 403)
	}
	res.Vlaid = true
	return nil
}

func (srv *UserService) publishEvent(reset *pb.PasswordReset) error {
	body, err := json.Marshal(reset)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Header: map[string]string{
			"email": reset.Email,
		},
		Body: body,
	}
	if err := srv.PubSub.Publish(passwordResetTopic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}
