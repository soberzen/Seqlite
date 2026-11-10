package model

import "time"

// DocumentPermissions 文档权限
type DocumentPermissions struct {
	ID          uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	Code        string    `gorm:"type:varchar(64);not null; uniqueIndex"`
	Name        string    `gorm:"type:varchar(128);not null"`
	Description string    `gorm:"type:text;not null;default:''"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

func (DocumentPermissions) TableName() string {
	return "document_permissions"
}

// SpacePermissions 空间权限
type SpacePermissions struct {
	ID          uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	Code        string    `gorm:"type:varchar(64);not null; uniqueIndex"`
	Name        string    `gorm:"type:varchar(128);not null"`
	Description string    `gorm:"type:text;not null;default:''"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

func (SpacePermissions) TableName() string {
	return "space_permissions"
}
