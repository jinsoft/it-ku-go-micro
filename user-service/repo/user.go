package repo

import (
	"github.com/jinsoft/it-ku/user-service/model"
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"gorm.io/gorm"
	"strconv"
)

type Repository interface {
	Create(user *model.User) error
	Get(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]*model.User, error)
}

type UserRepository struct {
	Db *gorm.DB
}

func (repo *UserRepository) Create(user *model.User) error {
	if err := repo.Db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) Get(id uint) (*model.User, error) {
	var user *pb.User
	sid := strconv.FormatInt(int64(id), 10)
	user.Id = sid
	if err := repo.Db.First(&user).Error; err != nil {
		return nil, err
	}
	userModel := &model.User{}
	muser, _ := userModel.ToORM(user)
	return muser, nil
}

func (repo *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &pb.User{}
	// 如果使用model.User{} 在查询的时候语句是  SELECT * FROM `users` WHERE email = 'admin@ainiok.com' AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
	// 改成不带有 deleted_at 字段的 pb.User{}
	if err := repo.Db.Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	userModel := &model.User{}
	muser, _ := userModel.ToORM(user)
	return muser, nil
}

func (repo *UserRepository) GetAll() ([]*model.User, error) {
	var users []*model.User
	if err := repo.Db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
