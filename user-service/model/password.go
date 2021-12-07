package model

import "time"

type PasswordReset struct {
	Email     string `gorm:"index"`
	Token     string `gorm:"not null"`
	CreatedAt time.Time
}

func (model *PasswordReset) ToORM() {

}
