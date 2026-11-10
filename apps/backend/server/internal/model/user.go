package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint           `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserName  string         `gorm:"type:varchar(32);not null"`
	Email     string         `gorm:"type:varchar(128);not null;index:idx_users_email,unique,where:deleted_at IS NULL"`
	Password  string         `gorm:"type:varchar(255);not null"`
	AvatarURL *string        `gorm:"type:varchar(255);"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index:idx_users_deleted_at"`
}

func (Users) TableName() string {
	return "users"
}
