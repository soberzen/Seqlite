package model

import (
	"time"

	"gorm.io/gorm"
)

type Comments struct {
	ID         uint           `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID     uint           `gorm:"type:bigint;not null;index:idx_comments_user_id"`
	DocumentID uint           `gorm:"type:bigint;not null;index:idx_comments_document_id"`
	ParentID   *uint          `gorm:"type:bigint;index:idx_comments_parent_id"` // 父评论 ID，自引用 Comments.ID
	Content    string         `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null;autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz;index:idx_comments_deleted_at"`
}

func (Comments) TableName() string {
	return "comments"
}
