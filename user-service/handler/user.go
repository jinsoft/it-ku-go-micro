package handler

import (
	"context"
	"github.com/jinsoft/it-ku/user-service/model"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"github.com/jinsoft/it-ku/user-service/repo"
	"github.com/jinsoft/it-ku/user-service/service"
	"github.com/micro/go-micro/v2/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

type UserService struct {
	Repo  repo.Repository
	Token service.Authable
}

func (srv *UserService) GetById(ctx context.Context, user *pb.User, response *pb.Response) error {
	panic("implement me")
}

func (srv *UserService) CreatePasswordReset(ctx context.Context, reset *pb.PasswordReset, response *pb.PasswordResetResponse) error {
	panic("implement me")
}

func (srv *UserService) ValidatePasswordResetToken(ctx context.Context, token *pb.Token, token2 *pb.Token) error {
	panic("implement me")
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

	if err != nil {
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
	user, _ := srv.Repo.Get(uint(id))
	if err := srv.Repo.Update(user); err != nil {
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
