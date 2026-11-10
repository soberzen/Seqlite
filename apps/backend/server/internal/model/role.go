package model

import "time"

type RoleCode string

const (
	RoleOwner     RoleCode = "owner"     // 所有者
	RoleAdmin     RoleCode = "admin"     // 管理员
	RoleEditor    RoleCode = "editor"    // 编辑者
	RoleViewer    RoleCode = "viewer"    // 查看者
	RoleCommenter RoleCode = "commenter" // 评论者
)

type Roles struct {
	ID          uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	Code        RoleCode  `gorm:"type:varchar(32);not null;uniqueIndex"`
	Name        string    `gorm:"type:varchar(64);not null"`
	Description string    `gorm:"type:text;not null;default:''"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}

func (Roles) TableName() string {
	return "roles"
}
