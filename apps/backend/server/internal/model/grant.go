package model

import "time"

type SpaceGrants struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID    uint      `gorm:"type:bigint;not null;uniqueIndex:uq_space_grants_user_space,priority:1"`
	SpaceID   uint      `gorm:"type:bigint;not null;uniqueIndex:uq_space_grants_user_space,priority:2;index:idx_space_grants_space_id"`
	RoleID    uint      `gorm:"type:bigint;not null;index:idx_space_grants_role_id"`
	GrantedBy uint      `gorm:"type:bigint;not null;index:idx_space_grants_granted_by"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}

func (SpaceGrants) TableName() string { return "space_grants" }

type DocumentGrants struct {
	ID         uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID     uint      `gorm:"type:bigint;not null;uniqueIndex:uq_document_grants_user_document,priority:1"`
	DocumentID uint      `gorm:"type:bigint;not null;uniqueIndex:uq_document_grants_user_document,priority:2;index:idx_document_grants_document_id"`
	RoleID     uint      `gorm:"type:bigint;not null;index:idx_document_grants_role_id"`
	GrantedBy  uint      `gorm:"type:bigint;not null;index:idx_document_grants_granted_by"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}

func (DocumentGrants) TableName() string { return "document_grants" }
