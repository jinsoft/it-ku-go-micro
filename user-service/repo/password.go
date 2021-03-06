package repo

import (
	pb "github.com/jinsoft/it-ku/user-service/proto/user"
	"gorm.io/gorm"
)

type PasswordResetInterface interface {
	Create(reset *pb.PasswordReset) error
	GetByToken(token string) (*pb.PasswordReset, error)
}

type PasswordResetRepository struct {
	Db *gorm.DB
}

func (repo *PasswordResetRepository) Create(reset *pb.PasswordReset) error {
	if err := repo.Db.Create(reset).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PasswordResetRepository) GetByToken(token string) (*pb.PasswordReset, error) {
	reset := &pb.PasswordReset{}
	if err := repo.Db.Where("token = ?", token).First(&reset).Error; err != nil {
		return nil, err
	}
	return reset, nil
}
