package model

import (
	"time"

	"github.com/google/uuid"
)

type DocumentStatus int8

const (
	DocumentStatusNormal  DocumentStatus = 0
	DocumentStatusTrashed DocumentStatus = 1
)

type Documents struct {
	ID      uint           `gorm:"type:bigint;primaryKey;autoIncrement"`
	SpaceID uint           `gorm:"type:bigint;not null;index:idx_documents_space_id"`
	Title   string         `gorm:"type:varchar(255);not null"`
	Emoji   *string        `gorm:"type:varchar(32)"`
	Status  DocumentStatus `gorm:"type:smallint;not null;default:0"`
	OwnerID uint           `gorm:"type:bigint;not null;index:idx_documents_owner_id"`
	// ExpireAt 仅当 Status=Trashed 时非 nil，表示回收站过期时间。
	ExpireAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz;not null;autoUpdateTime"`
	YjsRevision uint       `gorm:"column:yjs_revision;type:bigint;not null;default:0"`
	PublicID    uuid.UUID  `gorm:"type:uuid;not null;default:gen_random_uuid()"`
	// Owner     Users      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
}

func (Documents) TableName() string {
	return "documents"
}
