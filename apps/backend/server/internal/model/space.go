package model

import (
	"time"

	"github.com/google/uuid"
)

type SpaceStatus int8

const (
	SpaceStatusNormal  SpaceStatus = 0
	SpaceStatusTrashed SpaceStatus = 1
)

type Spaces struct {
	ID          uint        `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name        string      `gorm:"type:varchar(255);not null"`
	Emoji       *string     `gorm:"type:varchar(32)"`
	Description string      `gorm:"type:text;not null;default:''"`
	Status      SpaceStatus `gorm:"type:smallint;not null;default:0"`
	OwnerID     uint        `gorm:"type:bigint;not null;index:idx_spaces_owner_id"`
	PublicID    uuid.UUID   `gorm:"type:uuid;not null;default:gen_random_uuid()"`
	// ExpireAt 仅当 Status=Trashed 时非 nil，表示回收站过期时间。
	ExpireAt  *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;not null;autoUpdateTime"`
	// Owner     Users      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
}

func (Spaces) TableName() string {
	return "spaces"
}
