package model

import (
	"database/sql"
	"time"
)

type DeletedAt sql.NullTime

type BaseModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type SoftDelete struct {
	IsDeleted  uint8      `gorm:"column:is_deleted;default:0"`
	DeleteTime *time.Time `gorm:"column:delete_time;default:null"`
}
